package debuggercore

import (
	"encoding/binary"
	"github.com/kpli0rn/jdwpgo/protocol/basetypes"
	"github.com/kpli0rn/jdwpgo/protocol/vm"
)

// VMCommands expose the VM commands
type VMCommands interface {
	// Class
	AllClasses() (*vm.AllClassReply, error)
	// Thread ops
	AllThreads() (*vm.AllThreadsReply, error)
	TopLevelThreadGroups() (*vm.TopLevelThreadGroupsReply, error)
	StatusThread(threadId uint64) (*vm.ThreadStatusReply, error)

	// Bootstrap
	Version() (*vm.VersionReply, error)
	IDSizes() (*vm.IDSizesReply, error)
	Capabilities() (*vm.CapabilitiesReply, error)
	CapabilitiesNew() (*vm.CapabilitiesNewReply, error)
	//Control
	Suspend() error
	Resume() error
	HoldEvents() error
	ReleaseEvents() error
	Exit(int32) error

	// Methods
	AllMethods(refTypeId basetypes.JWDPRefTypeID) (*vm.AllMethodsReply, error)
	SendEventRequest(eventKind int8, threadId uint64) (*vm.EventRequestSetReply, error)
	ClearCommand(requestId int32) error

	// method invoke
	CreateString(command string) (*vm.CreateStringReply, error)
	// 无参
	InvokeStaticMethod(runtimeClasId basetypes.JWDPRefTypeID, threadId uint64, methodId uint64) (*vm.InvokeStaticMethodReply, error)
	// 有参
	InvokeMethod(objectID uint64, threadId uint64, runtimeClasId basetypes.JWDPRefTypeID, methodId uint64, commandID []byte) (*vm.InvokeMethodReply, error)
}

func (d *debuggercore) Version() (*vm.VersionReply, error) {
	var versionReply vm.VersionReply
	err := d.processCommand(vm.VersionCommand, nil, &versionReply, false)
	if err != nil {
		return nil, err
	}
	return &versionReply, nil
}

func (d *debuggercore) AllClasses() (*vm.AllClassReply, error) {
	var allclassesReply vm.AllClassReply
	err := d.processCommand(vm.AllClassesCommand, nil, &allclassesReply, false)
	if err != nil {
		return nil, err
	}
	return &allclassesReply, nil
}

func (d *debuggercore) AllThreads() (*vm.AllThreadsReply, error) {
	var allthreadsReply vm.AllThreadsReply
	err := d.processCommand(vm.AllThreadsCommand, nil, &allthreadsReply, false)
	if err != nil {
		return nil, err
	}
	return &allthreadsReply, nil
}

func (d *debuggercore) TopLevelThreadGroups() (*vm.TopLevelThreadGroupsReply, error) {
	var topLevelThreadGroupsReply vm.TopLevelThreadGroupsReply
	err := d.processCommand(vm.TopLevelThreadGroupsCommand, nil, &topLevelThreadGroupsReply, false)
	if err != nil {
		return nil, err
	}
	return &topLevelThreadGroupsReply, nil
}

func (d *debuggercore) IDSizes() (*vm.IDSizesReply, error) {
	var idsizesReply vm.IDSizesReply

	err := d.processCommand(vm.IDSizesCommand, nil, &idsizesReply, false)
	if err != nil {
		return nil, err
	}
	return &idsizesReply, nil
}

func (d *debuggercore) Capabilities() (*vm.CapabilitiesReply, error) {
	var capsReply vm.CapabilitiesReply
	err := d.processCommand(vm.CapabilitiesCommand, nil, &capsReply, false)
	if err != nil {
		return nil, err
	}
	return &capsReply, nil
}

func (d *debuggercore) CapabilitiesNew() (*vm.CapabilitiesNewReply, error) {
	var capsNewReply vm.CapabilitiesNewReply
	err := d.processCommand(vm.CapabilitiesNewCommand, nil, &capsNewReply, false)
	if err != nil {
		return nil, err
	}
	return &capsNewReply, nil
}

func (d *debuggercore) Suspend() error {
	err := d.processCommand(vm.SuspendCommand, nil, nil, false)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) Resume() error {
	err := d.processCommand(vm.ResumeCommand, nil, nil, false)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) HoldEvents() error {
	err := d.processCommand(vm.HoldEventsCommand, nil, nil, false)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) ReleaseEvents() error {
	err := d.processCommand(vm.ReleaseEventsCommand, nil, nil, false)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) Exit(code int32) error {
	exitCommandData := &vm.ExitCommandData{
		ExitCode: code,
	}
	err := d.processCommand(vm.ExitCommand, exitCommandData, nil, false)
	if err != nil {
		return err
	}
	return nil
}

// 需要创建 methods 然后发送并进行接受
func (d *debuggercore) AllMethods(refTypeId basetypes.JWDPRefTypeID) (*vm.AllMethodsReply, error) {
	var allmethodsReply vm.AllMethodsReply

	getCommandMethod := &vm.GetCommandMethod{
		RefType: refTypeId,
	}
	// 第一个是传统的模式，第二个是data的数据，第三个是处理返回的数据
	err := d.processCommand(vm.AllMethodsCommand, getCommandMethod, &allmethodsReply, false)
	if err != nil {
		return nil, err
	}
	return &allmethodsReply, nil
}

func (d *debuggercore) SendEventRequest(eventKind int8, threadId uint64) (*vm.EventRequestSetReply, error) {
	var eventRequestSetReply vm.EventRequestSetReply

	// 0f 01 一位
	setEvtRequest := &vm.SetEventRequest{
		EventKind:  eventKind, // 1位
		SuspendAll: 2,         // 1位
		Modifiers:  1,

		ModKind:  10, // 1位
		ThreadID: threadId,
		Size:     0,
		Depth:    0,
	}

	err := d.processCommand(vm.EventRequestCommand, setEvtRequest, &eventRequestSetReply, false)
	if err != nil {
		return nil, err
	}
	return &eventRequestSetReply, nil
}

// https://docs.oracle.com/en/java/javase/11/docs/specs/jdwp/jdwp-protocol.html#JDWP_EventRequest
func (d *debuggercore) ClearCommand(requestId int32) error {

	clearEventRequest := &vm.ClearEventRequest{
		EventKind: 1, // 这个eventkind 可以写死
		RequestID: requestId,
	}
	err := d.processCommand(vm.ClearEventCommand, clearEventRequest, nil, false)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) StatusThread(threadId uint64) (*vm.ThreadStatusReply, error) {
	var threadStatusReply vm.ThreadStatusReply

	threadStatusRequest := &vm.ThreadStatusRequest{
		ThreadID: threadId,
	}

	err := d.processCommand(vm.ThreadStatusCommand, threadStatusRequest, &threadStatusReply, false)
	if err != nil {
		return nil, err
	}
	return &threadStatusReply, err
}

// https://docs.oracle.com/javase/8/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_VirtualMachine_CreateString
func (d *debuggercore) CreateString(command string) (*vm.CreateStringReply, error) {
	var createStringReply vm.CreateStringReply
	lenCommand := make([]byte, 4)
	binary.BigEndian.PutUint32(lenCommand, uint32(len(command)))
	commandByte := append(lenCommand, []byte(command)...)
	//test, _ := hex.DecodeString("0000003d62617368202d63207b6563686f2c6233426c62694174595342445957786a64577868644739797d7c7b6261736536342c2d647d7c7b626173682c2d697d")
	err := d.processCommand(vm.CreateStringCommand, commandByte, &createStringReply, true)
	if err != nil {
		return nil, err
	}
	return &createStringReply, err
}

// https://docs.oracle.com/javase/8/docs/platform/jpda/jdwp/jdwp-protocol.html#JDWP_ClassType_InvokeMethod
func (d *debuggercore) InvokeStaticMethod(runtimeClasId basetypes.JWDPRefTypeID, threadId uint64, methodId uint64) (*vm.InvokeStaticMethodReply, error) {
	var invokeStaticMethodReply vm.InvokeStaticMethodReply

	invokeStaticMethodRequest := &vm.InvokeStaticMethodRequest{
		ClassID:  runtimeClasId,
		ThreadID: threadId,
		MethodID: methodId,
		ArgLen:   0,
		Options:  0,
	}

	err := d.processCommand(vm.InvokeStaticMethodCommand, invokeStaticMethodRequest, &invokeStaticMethodReply, false)
	if err != nil {
		return nil, err
	}
	return &invokeStaticMethodReply, err
}

func (d *debuggercore) InvokeMethod(objectID uint64, threadId uint64, runtimeClasId basetypes.JWDPRefTypeID, methodId uint64, commandID []byte) (*vm.InvokeMethodReply, error) {
	var invokeMethodReply vm.InvokeMethodReply

	invokeMethodRequest := &vm.InvokeMethodRequest{
		ObjectID: objectID,
		ThreadID: threadId,
		ClassID:  runtimeClasId,
		MethodID: methodId,
		ArgLen:   1,
		Arg:      commandID, // 这里应该要传递的是 id
		Options:  0,
	}
	err := d.processCommand(vm.InvokeMethodCommand, invokeMethodRequest, &invokeMethodReply, false)
	if err != nil {
		return nil, err
	}
	return &invokeMethodReply, err
}
