package debuggercore

import (
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
}

func (d *debuggercore) Version() (*vm.VersionReply, error) {
	var versionReply vm.VersionReply
	err := d.processCommand(vm.VersionCommand, nil, &versionReply)
	if err != nil {
		return nil, err
	}
	return &versionReply, nil
}

func (d *debuggercore) AllClasses() (*vm.AllClassReply, error) {
	var allclassesReply vm.AllClassReply
	err := d.processCommand(vm.AllClassesCommand, nil, &allclassesReply)
	if err != nil {
		return nil, err
	}
	return &allclassesReply, nil
}

func (d *debuggercore) AllThreads() (*vm.AllThreadsReply, error) {
	var allthreadsReply vm.AllThreadsReply
	err := d.processCommand(vm.AllThreadsCommand, nil, &allthreadsReply)
	if err != nil {
		return nil, err
	}
	return &allthreadsReply, nil
}

func (d *debuggercore) TopLevelThreadGroups() (*vm.TopLevelThreadGroupsReply, error) {
	var topLevelThreadGroupsReply vm.TopLevelThreadGroupsReply
	err := d.processCommand(vm.TopLevelThreadGroupsCommand, nil, &topLevelThreadGroupsReply)
	if err != nil {
		return nil, err
	}
	return &topLevelThreadGroupsReply, nil
}

func (d *debuggercore) IDSizes() (*vm.IDSizesReply, error) {
	var idsizesReply vm.IDSizesReply

	err := d.processCommand(vm.IDSizesCommand, nil, &idsizesReply)
	if err != nil {
		return nil, err
	}
	return &idsizesReply, nil
}

func (d *debuggercore) Capabilities() (*vm.CapabilitiesReply, error) {
	var capsReply vm.CapabilitiesReply
	err := d.processCommand(vm.CapabilitiesCommand, nil, &capsReply)
	if err != nil {
		return nil, err
	}
	return &capsReply, nil
}

func (d *debuggercore) CapabilitiesNew() (*vm.CapabilitiesNewReply, error) {
	var capsNewReply vm.CapabilitiesNewReply
	err := d.processCommand(vm.CapabilitiesNewCommand, nil, &capsNewReply)
	if err != nil {
		return nil, err
	}
	return &capsNewReply, nil
}

func (d *debuggercore) Suspend() error {
	err := d.processCommand(vm.SuspendCommand, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) Resume() error {
	err := d.processCommand(vm.ResumeCommand, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) HoldEvents() error {
	err := d.processCommand(vm.HoldEventsCommand, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) ReleaseEvents() error {
	err := d.processCommand(vm.ReleaseEventsCommand, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *debuggercore) Exit(code int32) error {
	exitCommandData := &vm.ExitCommandData{
		ExitCode: code,
	}
	err := d.processCommand(vm.ExitCommand, exitCommandData, nil)
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
	err := d.processCommand(vm.AllMethodsCommand, getCommandMethod, &allmethodsReply)
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

	err := d.processCommand(vm.EventRequestCommand, setEvtRequest, &eventRequestSetReply)
	if err != nil {
		return nil, err
	}
	return &eventRequestSetReply, nil
}

// https://docs.oracle.com/en/java/javase/11/docs/specs/jdwp/jdwp-protocol.html#JDWP_EventRequest
// eventkind 1 未测试 存在问题
func (d *debuggercore) ClearCommand(requestId int32) error {

	clearEventRequest := &vm.ClearEventRequest{
		EventKind: byte(1), // 这个eventkind 可以写死
		RequestID: requestId,
	}

	err := d.processCommand(vm.ClearEventCommand, clearEventRequest, nil)
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
	err := d.processCommand(vm.ThreadStatusCommand, threadStatusRequest, &threadStatusReply)
	if err != nil {
		return nil, err
	}
	return &threadStatusReply, err
}
