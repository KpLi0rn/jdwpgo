package vm

import (
	"fmt"
	"github.com/kpli0rn/jdwpgo/api/jdwp"
	"github.com/kpli0rn/jdwpgo/protocol/basetypes"
	"strings"
)

var AllMethodsCommand = jdwp.Command{Commandset: 2, Command: 5, HasCommandData: true, HasReplyData: true}

// https://docs.oracle.com/javase/7/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_ReferenceType_Methods
type AllMethodsReply struct {
	Declared int32
	Methods  []AllMethodsMethod `struct:"sizefrom=Declared"`
}

func (a *AllMethodsReply) String() string {
	var builder strings.Builder
	for _, method := range a.Methods {
		builder.WriteString(fmt.Sprintf("{%s}\n", method.String()))
	}
	return builder.String()
}

type AllMethodsMethod struct {
	MethodID  uint64
	Name      basetypes.JDWPString
	Signature basetypes.JDWPString
	ModBits   uint32
}

func (a *AllMethodsMethod) String() string {
	return fmt.Sprintf("MethodID: %v Name: %s Signature: %s ModBits: %v",
		a.MethodID,
		a.Name.String(),
		a.Signature.String(),
		a.ModBits,
	)
}
