package vm

import (
	"fmt"
	"github.com/jquirke/jdwpgo/api/jdwp"
)

var EventRequestCommand = jdwp.Command{Commandset: 15, Command: 1, HasCommandData: true, HasReplyData: true}

type EventRequestSetReply struct {
	RequestID int32
}

func (a *EventRequestSetReply) String() string {
	return fmt.Sprintf("RequestID: %v", a.RequestID)
}
