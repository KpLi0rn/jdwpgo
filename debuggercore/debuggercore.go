package debuggercore

import (
	"encoding/binary"
	"errors"
	"github.com/kpli0rn/jdwpgo/api/jdwp"
	"github.com/kpli0rn/jdwpgo/jdwpsession"
	"gopkg.in/restruct.v1"
)

// DebuggerCore represents an instance of the debugger core
type DebuggerCore interface {
	VMCommands() VMCommands
	ThreadCommands() ThreadCommands
}

type debuggercore struct {
	jdwpsession jdwpsession.Session
}

// NewFromJWDPSession creates a new instance of a debugger core
// attached to a JWDP session
func NewFromJWDPSession(session jdwpsession.Session) DebuggerCore {
	core := &debuggercore{
		jdwpsession: session,
	}

	return core
}

func (d *debuggercore) VMCommands() VMCommands {
	return d
}

func (d *debuggercore) ThreadCommands() ThreadCommands {
	return d
}

func (d *debuggercore) processCommand(cmd jdwp.Command, requestStruct interface{}, replyStruct interface{}, listen bool) error {
	commandPacket := &jdwpsession.CommandPacket{
		Commandset: cmd.Commandset,
		Command:    cmd.Command,
	}
	var err error
	if cmd.HasCommandData {
		commandPacket.Data, err = restruct.Pack(binary.BigEndian, requestStruct)
		if err != nil {
			return err
		}
	}
	// TODO implement timeout

	replyCh := d.jdwpsession.SendCommand(commandPacket)
	//select {
	//case reply, ok := <-replyCh:
	//	if !ok {
	//		return errors.New("Channel closed")
	//	}
	//	if cmd.HasReplyData {
	//		//_ = reply.Data
	//		err = restruct.Unpack(reply.Data, binary.BigEndian, replyStruct)
	//	}
	//	return err
	//
	//}

	if listen {
		return err
	}
	reply, ok := <-replyCh // 这里被阻塞了
	if !ok {
		return errors.New("Channel closed")
	}

	// TODO handle protocol returned err

	if cmd.HasReplyData {
		//_ = reply.Data
		err = restruct.Unpack(reply.Data, binary.BigEndian, replyStruct)
	}
	return err
}
