package vm

import (
	"fmt"
	"github.com/kpli0rn/jdwpgo/api/jdwp"
)

// 静态方法的调用
var InvokeStaticMethodCommand = jdwp.Command{Commandset: 3, Command: 3, HasCommandData: true, HasReplyData: true}

type InvokeStaticMethodReply struct {
	Tag       int8
	ContextID uint64
}

func (t *InvokeStaticMethodReply) String() string {
	return fmt.Sprintf("Tag: %v ContextID: %v", t.Tag, t.ContextID)
}

// 普通方法调用
var InvokeMethodCommand = jdwp.Command{Commandset: 9, Command: 6, HasCommandData: true, HasReplyData: true}

type InvokeMethodReply struct {
	Tag       int8
	ContextID uint64
}

func (t *InvokeMethodReply) String() string {
	return fmt.Sprintf("Tag: %v ContextID: %v", t.Tag, t.ContextID)
}
