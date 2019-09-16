package vm

import (
	"fmt"
	"strings"

	"github.com/jquirke/jdwpgo/protocol/common"

	"github.com/jquirke/jdwpgo/api/jdwp"
)

// AllThreadsCommand represents the AllThreads command
var AllThreadsCommand = jdwp.Command{Commandset: 1, Command: 4, HasReplyData: true}

// AllThreadsReply represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_AllThreads
type AllThreadsReply struct {
	NumThreads int32
	Threads    []common.ThreadID `struct:"sizefrom=NumThreads"`
}

func (a *AllThreadsReply) String() string {
	var builder strings.Builder
	for _, threadID := range a.Threads {
		builder.WriteString(fmt.Sprintf("%s\n", threadID.String()))
	}
	return builder.String()
}

// TopLevelThreadGroupsCommand represents the TopLevelThreadsGroups command
var TopLevelThreadGroupsCommand = jdwp.Command{Commandset: 1, Command: 5, HasReplyData: true}

// TopLevelThreadGroupsReply represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_TopLevelThreadGroups
type TopLevelThreadGroupsReply struct {
	NumThreadGroups int32
	ThreadGroups    []common.ThreadGroupID `struct:"sizefrom=NumThreadGroups"`
}

func (t *TopLevelThreadGroupsReply) String() string {
	var builder strings.Builder
	for _, threadGroupID := range t.ThreadGroups {
		builder.WriteString(fmt.Sprintf("%s\n", threadGroupID.String()))
	}
	return builder.String()
}
