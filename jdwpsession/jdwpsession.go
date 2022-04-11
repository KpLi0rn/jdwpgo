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
const defaultReadDeadlineMillis = 2000
const defaultWriteDeadlineMillis = 2000

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
	SendCommand(*CommandPacket) <-chan *ReplyPacket
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
	id            uint32
	flags         byte
	commandPacket *CommandPacket
	replyPacket   *ReplyPacket
}

func (w *WrappedPacket) isCommandPacket() bool {
	return w.commandPacket != nil
}

func (w *WrappedPacket) String() string {
	if w.commandPacket != nil {
		return fmt.Sprintf("{id=%v flags=%x commandpacket=%v", w.id, w.flags, w.commandPacket)
	}
	return fmt.Sprintf("{id=%v flags=%x replypacket=%v", w.id, w.flags, w.replyPacket)
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
		requestPendingQueue: make(chan *request, 10),
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
	go s.rxLoop()
	go s.txLoop()
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

func (s *session) rxLoop() {
	for atomic.LoadInt32(&s.state) == sessionOpen {
		err := s.dispatchInboundPacket()
		if err != nil {
			s.setErrorState(err)
			break
		}
	}

	close(s.jvmCommandPackets)
	for _, request := range s.requestPending {
		close(request.replyCh)
	}
}

func (s *session) setErrorState(err error) {
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()
	fmt.Printf("closing session due to error: %v\n", err)
	if s.state == sessionOpen {
		s.state = sessionFailed
	}
}

func (s *session) txLoop() {
	// TODO need exit from here
	for request := range s.requestPendingQueue {
		err := s.writePacket(request)
		if err != nil {
			s.setErrorState(err)
			break
		}
	}
}

func (s *session) writePacket(request *request) error {
	s.conn.SetWriteDeadline(time.Now().Add(defaultWriteDeadlineMillis * time.Millisecond))
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

func (s *session) dispatchInboundPacket() error {
	wrappedPacket, err := s.readPacket()
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()

	if err != nil {
		return err
	}
	if wrappedPacket.isCommandPacket() {
		s.jvmCommandPackets <- wrappedPacket.commandPacket
	} else {
		request, ok := s.requestPending[wrappedPacket.id]
		if !ok {
			fmt.Printf("warn: got unexpected reply for id: %v", wrappedPacket.id)
		} else {
			request.replyCh <- wrappedPacket.replyPacket
			//close(request.replyCh) //TODO turn back on
		}
	}
	return nil
}
func (s *session) readPacket() (*WrappedPacket, error) {
	var wrappedPacket WrappedPacket
	s.conn.SetReadDeadline(time.Time{})
	var size uint32
	err := binary.Read(s.conn, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	s.conn.SetReadDeadline(time.Now().Add(defaultReadDeadlineMillis * time.Millisecond))
	if size < headerBytes {
		return nil, fmt.Errorf("packet too small: %v", size)
	}
	dataSize := size - headerBytes
	err = binary.Read(s.conn, binary.BigEndian, &wrappedPacket.id)
	if err != nil {
		return nil, err
	}
	err = binary.Read(s.conn, binary.BigEndian, &wrappedPacket.flags)
	if err != nil {
		return nil, err
	}

	var dataSlice *[]byte
	if wrappedPacket.flags&flagsReplyPacket == flagsReplyPacket {
		var replyPacket ReplyPacket
		wrappedPacket.replyPacket = &replyPacket
		err = binary.Read(s.conn, binary.BigEndian, &replyPacket.Errorcode)
		if err != nil {
			return nil, err
		}
		dataSlice = &replyPacket.Data
	} else {
		var commandPacket CommandPacket
		wrappedPacket.commandPacket = &commandPacket
		err = binary.Read(s.conn, binary.BigEndian, &commandPacket.Commandset)
		if err != nil {
			return nil, err
		}
		err = binary.Read(s.conn, binary.BigEndian, &commandPacket.Command)
		if err != nil {
			return nil, err
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

func (s *session) SendCommand(commandPacket *CommandPacket) <-chan *ReplyPacket {
	fmt.Println(commandPacket)
	sendid := atomic.AddUint32(&s.sequence, 1)
	request := request{
		id:            sendid,
		replyCh:       make(chan *ReplyPacket, 1),
		commandPacket: commandPacket,
	}
	s.sessionMutex.Lock()
	s.requestPending[sendid] = &request
	s.sessionMutex.Unlock()
	// the transmission MUST occur after

	s.requestPendingQueue <- &request

	return request.replyCh
}
