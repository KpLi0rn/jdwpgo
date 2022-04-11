package vm

import (
	"github.com/jquirke/jdwpgo/api/jdwp"
	"github.com/jquirke/jdwpgo/protocol/basetypes"
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
}

// HoldEventsCommand represents the hold events command
var HoldEventsCommand = jdwp.Command{Commandset: 1, Command: 15}

// ReleaseEventsCommand represents the hold events command
var ReleaseEventsCommand = jdwp.Command{Commandset: 1, Command: 16}
