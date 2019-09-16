package thread

import (
	"github.com/jquirke/jdwpgo/api/jdwp"
	"github.com/jquirke/jdwpgo/protocol/basetypes"
	"github.com/jquirke/jdwpgo/protocol/common"
)

// NameCommand represents the name command
var NameCommand = jdwp.Command{Commandset: 11, Command: 1, HasCommandData: true, HasReplyData: true}

// NameCommandData represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_ThreadReference_Name
type NameCommandData struct {
	ThreadID common.ThreadID
}

// NameReply represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_ThreadReference_Name
type NameReply struct {
	ThreadName basetypes.JDWPString
}
