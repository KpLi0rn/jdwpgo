package vm

import (
	"fmt"

	"github.com/kpli0rn/jdwpgo/api/jdwp"
	"github.com/kpli0rn/jdwpgo/protocol/basetypes"
)

// VersionCommand represents the version command
var VersionCommand = jdwp.Command{Commandset: 1, Command: 1, HasReplyData: true}

// VersionReply represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_Version
type VersionReply struct {
	Description basetypes.JDWPString
	JwdpMajor   int32
	JwdpMinor   int32
	VMVersion   basetypes.JDWPString
	VMName      basetypes.JDWPString
}

func (v *VersionReply) String() string {
	return fmt.Sprintf("Description: %s\nJWDPMajor: %v\nJWDPMinor: %v\nVMVersion: %s\nVMName: %s\n",
		v.Description.String(),
		v.JwdpMajor,
		v.JwdpMinor,
		v.VMVersion.String(),
		v.VMName.String())
}
