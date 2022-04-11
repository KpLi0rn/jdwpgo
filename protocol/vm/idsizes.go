package vm

import (
	"fmt"

	"github.com/kpli0rn/jdwpgo/api/jdwp"
)

// IDSizesCommand represents the IDSizes command
var IDSizesCommand = jdwp.Command{Commandset: 1, Command: 7, HasReplyData: true}

// IDSizesReply represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_IDSizes
type IDSizesReply struct {
	FieldIDSize         int32
	MethodIDSize        int32
	ObjectIDSize        int32
	ReferenceTypeIDSize int32
	FrameIDSize         int32
}

func (i *IDSizesReply) String() string {
	return fmt.Sprintf("FieldIDSize: %v\nMethodIDSize: %v\nObjectIDSize: %v\nReferenceTypeIDSize: %v\nFrameIDSize: %v\n",
		i.FieldIDSize,
		i.MethodIDSize,
		i.ObjectIDSize,
		i.ReferenceTypeIDSize,
		i.FrameIDSize)
}
