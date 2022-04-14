package vm

import (
	"github.com/kpli0rn/jdwpgo/api/jdwp"
	"github.com/kpli0rn/jdwpgo/protocol/basetypes"
)

// SuspendCommand represents the suspendcommand
var SuspendCommand = jdwp.Command{Commandset: 1, Command: 8}

// ResumeCommand represents the resume command
var ResumeCommand = jdwp.Command{Commandset: 1, Command: 9}

// ExitCommand represents the hold events command
var ExitCommand = jdwp.Command{Commandset: 1, Command: 10, HasCommandData: true}

// ExitCommandData represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_Exit
type ExitCommandData struct {
	ExitCode int32
}

// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_ReferenceType_Methods
type GetCommandMethod struct {
	RefType basetypes.JWDPRefTypeID
}

type SetEventRequest struct {
	EventKind  int8
	SuspendAll int8
	Modifiers  int32
	ModKind    int8
	ThreadID   uint64 // uint64
	Size       int32
	Depth      int32
}

//type EventRequest struct {
//	ModKind  int32
//	ThreadID uint64 // uint64
//	Size     int32
//	Depth    int32
//}

// HoldEventsCommand represents the hold events command
var HoldEventsCommand = jdwp.Command{Commandset: 1, Command: 15}

// ReleaseEventsCommand represents the hold events command
var ReleaseEventsCommand = jdwp.Command{Commandset: 1, Command: 16}

var ClearEventCommand = jdwp.Command{Commandset: 15, Command: 2, HasCommandData: true, HasReplyData: false}

type ClearEventRequest struct {
	EventKind byte
	RequestID int32
}

type ThreadStatusRequest struct {
	ThreadID uint64
}

//type CreateStringRequest struct {
//	Command []byte
//}

type InvokeStaticMethodRequest struct {
	ClassID  basetypes.JWDPRefTypeID
	ThreadID uint64
	MethodID uint64
	ArgLen   int32
	//Arg      []string
	Options int32
}

type InvokeMethodRequest struct {
	ObjectID uint64
	ThreadID uint64
	ClassID  basetypes.JWDPRefTypeID
	MethodID uint64
	ArgLen   uint32
	Arg      []byte
	Options  int32
}
