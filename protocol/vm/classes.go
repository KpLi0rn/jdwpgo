package vm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kpli0rn/jdwpgo/api/jdwp"
	"github.com/kpli0rn/jdwpgo/protocol/basetypes"
)

// AllClassesCommand represents the all classes command
var AllClassesCommand = jdwp.Command{Commandset: 1, Command: 3, HasReplyData: true}

// AllClassReply represents
// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_AllClasses
type AllClassReply struct {
	NumClasses int32
	Classes    []AllClassClass `struct:"sizefrom=NumClasses"`
}

func (a *AllClassReply) String() string {
	var builder strings.Builder
	for _, class := range a.Classes {
		builder.WriteString(fmt.Sprintf("{%s}\n", class.String()))
	}
	return builder.String()
}

// AllClassClass represents a single class in AllClassReply
type AllClassClass struct {
	RefTypeTag      basetypes.JWDPTypeTag
	ReferenceTypeID basetypes.JWDPRefTypeID
	Signature       basetypes.JDWPString
	Status          AllClassClassStatus //TODO ENUM
}

func (a *AllClassClass) String() string {
	return fmt.Sprintf("RefTypeTag: %v ReferenceTypeID: %s Signature: %s Status: %v",
		a.RefTypeTag.String(),
		a.ReferenceTypeID.String(),
		a.Signature.String(),
		a.Status.String(),
	)
}

// check isempty
func (a *AllClassClass) IsEmpty() bool {
	return reflect.DeepEqual(a, AllClassClass{})
}

// AllClassClassStatus represents a class's status
type AllClassClassStatus int32

const (
	// AllClassClassStatusVerified - class verified
	AllClassClassStatusVerified = 1
	// AllClassClassStatusPrepared - class prepared
	AllClassClassStatusPrepared = 2
	// AllClassClassStatusInitialized - class initialized
	AllClassClassStatusInitialized = 4
	// AllClassClassStatusError - class error
	AllClassClassStatusError = 8
)

func (a AllClassClassStatus) String() string {
	labels := []struct {
		mask  int32
		label string
	}{
		{AllClassClassStatusVerified, "Verified"},
		{AllClassClassStatusPrepared, "Prepared"},
		{AllClassClassStatusInitialized, "Initialized"},
		{AllClassClassStatusError, "Error"},
	}
	var builder strings.Builder
	builder.WriteString("{")
	for idx, label := range labels {
		if (int32)(a)&label.mask != 0 {
			if idx != 0 {
				builder.WriteString("|")
			}
			builder.WriteString(label.label)
		}
	}
	builder.WriteString("}")
	return builder.String()
}
