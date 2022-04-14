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

func (d *debuggercore) processCommand(cmd jdwp.Command, requestStruct interface{}, replyStruct interface{}, custom bool) error {
	commandPacket := &jdwpsession.CommandPacket{
		Commandset: cmd.Commandset,
		Command:    cmd.Command,
	}
	var err error

	// 可自己控制输入
	if custom {
		commandPacket.Data = requestStruct.([]byte)
	}
	if cmd.HasCommandData {
		commandPacket.Data, err = restruct.Pack(binary.BigEndian, requestStruct)
		if err != nil {
			return err
		}
	}

	wrapPacket := d.jdwpsession.SendCommand(commandPacket)
	if wrapPacket == nil {
		return errors.New("reply is nil")
	}
	if cmd.HasReplyData {
		// 4c000000000000117c4c0000000000000000
		err = restruct.Unpack(wrapPacket.ReplyPacket.Data, binary.BigEndian, replyStruct)
	}
	return err
}
