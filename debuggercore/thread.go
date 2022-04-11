package debuggercore

import (
	"github.com/kpli0rn/jdwpgo/protocol/basetypes"
	"github.com/kpli0rn/jdwpgo/protocol/common"
	"github.com/kpli0rn/jdwpgo/protocol/thread"
)

// ThreadCommands expose the ThreadCommands commands
type ThreadCommands interface {
	// Basics
	Name(common.ThreadID) (basetypes.JDWPString, error)
}

func (d *debuggercore) Name(threadID common.ThreadID) (basetypes.JDWPString, error) {
	nameCommandData := &thread.NameCommandData{
		ThreadID: threadID,
	}
	var nameReply thread.NameReply
	err := d.processCommand(thread.NameCommand, nameCommandData, &nameReply)
	if err != nil {
		return basetypes.EmptyJWDPString(), err
	}
	return nameReply.ThreadName, nil
}
