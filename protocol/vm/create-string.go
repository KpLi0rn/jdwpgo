package vm

import (
	"fmt"
	"github.com/kpli0rn/jdwpgo/api/jdwp"
	"github.com/kpli0rn/jdwpgo/protocol/basetypes"
)

var CreateStringCommand = jdwp.Command{Commandset: 1, Command: 11, HasCommandData: false, HasReplyData: true}

type CreateStringReply struct {
	StringObject basetypes.JWDPObjectID
}

func (t *CreateStringReply) String() string {
	return fmt.Sprintf("StringObject: %v", t.StringObject)
}
