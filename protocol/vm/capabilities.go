package vm

import (
	"fmt"
	"strings"

	"github.com/kpli0rn/jdwpgo/api/jdwp"
)

// CapabilitiesCommand represents the capabilities command
var CapabilitiesCommand = jdwp.Command{Commandset: 1, Command: 12, HasReplyData: true}

// CapabilitiesReply represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_Capabilities
type CapabilitiesReply struct {
	CanWatchFieldModification     bool
	CanWatchFieldAccess           bool
	CanGetBytecodes               bool
	CanGetSyntheticAttribute      bool
	CanGetOwnedMonitorInfo        bool
	CanGetCurrentContendedMonitor bool
	CanGetMonitorInfo             bool
}

func (c *CapabilitiesReply) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("CanWatchFieldModification: %v\n", c.CanWatchFieldModification))
	builder.WriteString(fmt.Sprintf("CanWatchFieldAccess: %v\n", c.CanWatchFieldAccess))
	builder.WriteString(fmt.Sprintf("CanGetBytecodes: %v\n", c.CanGetBytecodes))
	builder.WriteString(fmt.Sprintf("CanGetSyntheticAttribute: %v\n", c.CanGetSyntheticAttribute))
	builder.WriteString(fmt.Sprintf("CanGetOwnedMonitorInfo: %v\n", c.CanGetOwnedMonitorInfo))
	builder.WriteString(fmt.Sprintf("CanGetCurrentContendedMonitor: %v\n", c.CanGetCurrentContendedMonitor))
	builder.WriteString(fmt.Sprintf("CanGetMonitorInfo: %v\n", c.CanGetMonitorInfo))

	return builder.String()
}

// CapabilitiesNewCommand represents the capabilities command
var CapabilitiesNewCommand = jdwp.Command{Commandset: 1, Command: 17, HasReplyData: true}

// CapabilitiesNewReply represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_CapabilitiesNew
type CapabilitiesNewReply struct {
	CanWatchFieldModification        bool
	CanWatchFieldAccess              bool
	CanGetBytecodes                  bool
	CanGetSyntheticAttribute         bool
	CanGetOwnedMonitorInfo           bool
	CanGetCurrentContendedMonitor    bool
	CanGetMonitorInfo                bool
	CanRedefineClasses               bool
	CanAddMethod                     bool
	CanUnrestrictedlyRedefineClasses bool
	CanPopFrames                     bool
	CanUseInstanceFilters            bool
	CanGetSourceDebugExtension       bool
	CanRequestVMDeathEvent           bool
	CanSetDefaultStratum             bool
	CanGetInstanceInfo               bool
	CanRequestMonitorEvents          bool
	CanGetMonitorFrameInfo           bool
	CanUseSourceNameFilters          bool
	CanGetConstantPool               bool
	CanForceEarlyReturn              bool
	Reserved22                       bool
	Reserved23                       bool
	Reserved24                       bool
	Reserved25                       bool
	Reserved26                       bool
	Reserved27                       bool
	Reserved28                       bool
	Reserved29                       bool
	Reserved30                       bool
	Reserved31                       bool
	Reserved32                       bool
}

func (c *CapabilitiesNewReply) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("CanWatchFieldModification: %v\n", c.CanWatchFieldModification))
	builder.WriteString(fmt.Sprintf("CanWatchFieldAccess: %v\n", c.CanWatchFieldAccess))
	builder.WriteString(fmt.Sprintf("CanGetBytecodes: %v\n", c.CanGetBytecodes))
	builder.WriteString(fmt.Sprintf("CanGetSyntheticAttribute: %v\n", c.CanGetSyntheticAttribute))
	builder.WriteString(fmt.Sprintf("CanGetOwnedMonitorInfo: %v\n", c.CanGetOwnedMonitorInfo))
	builder.WriteString(fmt.Sprintf("CanGetCurrentContendedMonitor: %v\n", c.CanGetCurrentContendedMonitor))
	builder.WriteString(fmt.Sprintf("CanGetMonitorInfo: %v\n", c.CanGetMonitorInfo))
	builder.WriteString(fmt.Sprintf("CanRedefineClasses: %v\n", c.CanRedefineClasses))
	builder.WriteString(fmt.Sprintf("CanAddMethod: %v\n", c.CanAddMethod))
	builder.WriteString(fmt.Sprintf("CanUnrestrictedlyRedefineClasses: %v\n", c.CanUnrestrictedlyRedefineClasses))
	builder.WriteString(fmt.Sprintf("CanPopFrames: %v\n", c.CanPopFrames))
	builder.WriteString(fmt.Sprintf("CanUseInstanceFilters: %v\n", c.CanUseInstanceFilters))
	builder.WriteString(fmt.Sprintf("CanGetSourceDebugExtension: %v\n", c.CanGetSourceDebugExtension))
	builder.WriteString(fmt.Sprintf("CanRequestVMDeathEvent: %v\n", c.CanRequestVMDeathEvent))
	builder.WriteString(fmt.Sprintf("CanSetDefaultStratum: %v\n", c.CanSetDefaultStratum))
	builder.WriteString(fmt.Sprintf("CanGetInstanceInfo: %v\n", c.CanGetInstanceInfo))
	builder.WriteString(fmt.Sprintf("CanRequestMonitorEvents: %v\n", c.CanRequestMonitorEvents))
	builder.WriteString(fmt.Sprintf("CanGetMonitorFrameInfo: %v\n", c.CanGetMonitorFrameInfo))
	builder.WriteString(fmt.Sprintf("CanUseSourceNameFilters: %v\n", c.CanUseSourceNameFilters))
	builder.WriteString(fmt.Sprintf("CanGetConstantPool: %v\n", c.CanGetConstantPool))
	builder.WriteString(fmt.Sprintf("CanForceEarlyReturn: %v\n", c.CanForceEarlyReturn))
	builder.WriteString(fmt.Sprintf("Reserved22: %v\n", c.Reserved22))
	builder.WriteString(fmt.Sprintf("Reserved23: %v\n", c.Reserved23))
	builder.WriteString(fmt.Sprintf("Reserved24: %v\n", c.Reserved24))
	builder.WriteString(fmt.Sprintf("Reserved25: %v\n", c.Reserved25))
	builder.WriteString(fmt.Sprintf("Reserved26: %v\n", c.Reserved26))
	builder.WriteString(fmt.Sprintf("Reserved27: %v\n", c.Reserved27))
	builder.WriteString(fmt.Sprintf("Reserved28: %v\n", c.Reserved28))
	builder.WriteString(fmt.Sprintf("Reserved29: %v\n", c.Reserved29))
	builder.WriteString(fmt.Sprintf("Reserved30: %v\n", c.Reserved30))
	builder.WriteString(fmt.Sprintf("Reserved31: %v\n", c.Reserved31))
	builder.WriteString(fmt.Sprintf("Reserved32: %v\n", c.Reserved32))

	return builder.String()
}
