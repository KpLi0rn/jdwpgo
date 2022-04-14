package jdwpsession

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const defaultPacketQueueLength = 50
const defaultReadDeadlineMillis = 5000
const defaultWriteDeadlineMillis = 5000

const headerBytes = 11
const handshakebytes = "JDWP-Handshake"
const flagsReplyPacket = 0x80

const (
	sessionClosed = iota
	sessionHandshake
	sessionOpen
	sessionFailed
)

// Session implements the low level JWDP session abstraction
// it is thread safe and supports concurrent in flight
// requests/responses
type Session interface {
	Start() error
	Stop() error
	JvmCommandPacketChannel() <-chan *CommandPacket
	//SendCommand(*CommandPacket) <-chan *ReplyPacket
	SendCommand(*CommandPacket) *WrappedPacket
	ReadPacket() (*WrappedPacket, error)
	WritePacket(*request) error

	DispatchInboundPacket() error
}

type session struct {
	conn              net.Conn
	jvmCommandPackets chan *CommandPacket
	sessionMutex      sync.Mutex
	// mutex protected
	requestPending      map[uint32]*request
	requestPendingQueue chan *request
	state               int32
	sequence            uint32
}

type request struct {
	id            uint32
	replyCh       chan *ReplyPacket
	commandPacket *CommandPacket
}

// WrappedPacket represents a command or reply packet
type WrappedPacket struct {
	Id            uint32
	Flags         byte
	CommandPacket *CommandPacket
	ReplyPacket   *ReplyPacket
}

func (w *WrappedPacket) isCommandPacket() bool {
	return w.CommandPacket != nil
}

func (w *WrappedPacket) String() string {
	if w.CommandPacket != nil {
		return fmt.Sprintf("{id=%v flags=%x commandpacket=%v", w.Id, w.Flags, w.CommandPacket)
	}
	return fmt.Sprintf("{id=%v flags=%x replypacket=%v", w.Id, w.Flags, w.ReplyPacket)
}

// CommandPacket represents a command packet
type CommandPacket struct {
	Commandset byte
	Command    byte
	Data       []byte
}

func (c *CommandPacket) String() string {
	return fmt.Sprintf("{commandset=%v[TODO] command=%v[TODO] length=%v",
		c.Commandset, c.Command, len(c.Data))
}

// ReplyPacket represents a reply packet
type ReplyPacket struct {
	Errorcode uint16
	Data      []byte
}

func (r *ReplyPacket) String() string {
	return fmt.Sprintf("{errorcode=%v length=%v",
		r.Errorcode, len(r.Data))
}

// New creates a new JWDP session
func New(conn net.Conn) Session {
	return &session{
		conn:                conn,
		requestPending:      make(map[uint32]*request),
		requestPendingQueue: make(chan *request, 10), // 请求池 容量 10
	}
}

func (s *session) Start() error {
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()

	if s.state != sessionClosed {
		return errors.New("session not in closed state")
	}
	s.state = sessionHandshake
	if err := s.writeHandshakeFrame(); err != nil {
		s.state = sessionFailed
		return err
	}
	if err := s.readAndCheckHandshakeFrame(); err != nil {
		s.state = sessionFailed
		return err
	}
	s.jvmCommandPackets = make(chan *CommandPacket, defaultPacketQueueLength)
	s.state = sessionOpen
	return nil
}

func (s *session) writeHandshakeFrame() error {
	s.conn.SetWriteDeadline(time.Now().Add(defaultWriteDeadlineMillis * time.Millisecond))
	_, err := s.conn.Write([]byte(handshakebytes))
	return err
}

func (s *session) readAndCheckHandshakeFrame() error {
	s.conn.SetReadDeadline(time.Now().Add(defaultReadDeadlineMillis * time.Millisecond))
	buf := make([]byte, len(handshakebytes))
	_, err := io.ReadFull(s.conn, buf)
	return err
}

func (s *session) setErrorState(err error) {
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()
	fmt.Printf("closing session due to error: %v\n", err)
	if s.state == sessionOpen {
		s.state = sessionFailed
	}
}

func (s *session) WritePacket(request *request) error {
	s.conn.SetWriteDeadline(time.Now().Add(defaultWriteDeadlineMillis * time.Millisecond))
	//s.conn.SetWriteDeadline(time.Now().Add(100 * time.Second))

	var totalsize = 11 + (uint32)(len(request.commandPacket.Data))
	err := binary.Write(s.conn, binary.BigEndian, totalsize)
	if err != nil {
		return err
	}
	err = binary.Write(s.conn, binary.BigEndian, request.id)
	if err != nil {
		return err
	}
	err = binary.Write(s.conn, binary.BigEndian, (byte)(0))
	if err != nil {
		return err
	}
	err = binary.Write(s.conn, binary.BigEndian, request.commandPacket.Commandset)
	if err != nil {
		return err
	}
	err = binary.Write(s.conn, binary.BigEndian, request.commandPacket.Command)
	if err != nil {
		return err
	}
	n, err := s.conn.Write(request.commandPacket.Data)
	if err != nil {
		return err
	}
	if n != len(request.commandPacket.Data) {
		return fmt.Errorf("did not write all bytes, got %v expect %v",
			n, len(request.commandPacket.Data))
	}
	return nil
}

//
func (s *session) DispatchInboundPacket() error {
	wrappedPacket, err := s.ReadPacket()
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()
	if err != nil {
		return err
	}
	if wrappedPacket.isCommandPacket() {
		s.jvmCommandPackets <- wrappedPacket.CommandPacket
	} else {
		request, ok := s.requestPending[wrappedPacket.Id]
		if !ok {
			fmt.Printf("warn: got unexpected reply for id: %v", wrappedPacket.Id)
		} else {
			request.replyCh <- wrappedPacket.ReplyPacket
			close(request.replyCh) //TODO turn back on
		}
	}
	return nil
}

func (s *session) ReadPacket() (*WrappedPacket, error) {

	var wrappedPacket WrappedPacket
	s.conn.SetReadDeadline(time.Now().Add(100 * time.Second))
	var size uint32
	err := binary.Read(s.conn, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	if size < headerBytes {
		return nil, fmt.Errorf("packet too small: %v", size)
	}
	dataSize := size - headerBytes
	err = binary.Read(s.conn, binary.BigEndian, &wrappedPacket.Id)
	if err != nil {
		return nil, errors.New("2")
	}
	err = binary.Read(s.conn, binary.BigEndian, &wrappedPacket.Flags)
	if err != nil {
		return nil, errors.New("3")
	}
	var dataSlice *[]byte
	if wrappedPacket.Flags&flagsReplyPacket == flagsReplyPacket {
		var replyPacket ReplyPacket
		wrappedPacket.ReplyPacket = &replyPacket
		err = binary.Read(s.conn, binary.BigEndian, &replyPacket.Errorcode)
		if err != nil {
			return nil, errors.New("4")
		}
		dataSlice = &replyPacket.Data
	} else { // 这里是处理请求
		var commandPacket CommandPacket
		wrappedPacket.CommandPacket = &commandPacket
		err = binary.Read(s.conn, binary.BigEndian, &commandPacket.Commandset)
		if err != nil {
			return nil, errors.New("5")
		}
		err = binary.Read(s.conn, binary.BigEndian, &commandPacket.Command)
		if err != nil {
			return nil, errors.New("6")
		}
		dataSlice = &commandPacket.Data
	}

	*dataSlice = make([]byte, dataSize)
	_, err = io.ReadFull(s.conn, *dataSlice)
	if err != nil {
		return nil, err
	}
	return &wrappedPacket, nil
}

func (s *session) Stop() error {
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()

	if s.state == sessionOpen {
		s.state = sessionClosed
	} else {
		return fmt.Errorf("session not open: %v", s.state)
	}
	return nil
}

func (s *session) JvmCommandPacketChannel() <-chan *CommandPacket {
	return s.jvmCommandPackets
}

func (s *session) SendCommand(commandPacket *CommandPacket) *WrappedPacket {
	sendid := atomic.AddUint32(&s.sequence, 1)
	request := request{
		id:            sendid,
		replyCh:       make(chan *ReplyPacket, 1),
		commandPacket: commandPacket,
	}
	// 对数据进行发送
	s.WritePacket(&request)
	wrappedPacket, _ := s.ReadPacket()
	//fmt.Println(wrappedPacket.ReplyPacket)
	s.sessionMutex.Lock()
	s.requestPending[sendid] = &request
	s.sessionMutex.Unlock()

	return wrappedPacket
}
