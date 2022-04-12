package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/kpli0rn/jdwpgo/debuggercore"
	"github.com/kpli0rn/jdwpgo/jdwpsession"
	"github.com/kpli0rn/jdwpgo/protocol/vm"
	"log"
	"net"
)

func getMethodByName(methods *vm.AllMethodsReply, name string) *vm.AllMethodsMethod {
	for _, method := range methods.Methods {
		if method.Name.String() == name {
			return &method
		}
	}
	return nil
}

// 40640200000001010000000200000000000011770100000000000003e100006000028a03e80000000000000031
// size: 50 针对这个进行解析，提取出来 1144
// 000000070040640200000001010000000200000000000011440100000000000003d50000600002e415a80000000000000031
// size: 54
// 000000040040640200000001010000000200000000000011420100000000000003d30000600003dca6680000000000000031
func parseEvent(buf []byte, eventId int32, idsize *vm.IDSizesReply) (int32, uint64) {
	raw := buf[7:]
	rId := int32(binary.BigEndian.Uint32(raw[6:10]))
	if rId != eventId {
		return 0, 0
	}
	rawtId := raw[10 : 10+idsize.ObjectIDSize]
	tId := binary.BigEndian.Uint64(rawtId)
	return rId, tId
}

func runtimeExec() {

}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
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
	idSizes, err := debuggerCore.VMCommands().IDSizes() // 这里的 IDsize 没有用起来
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
	log.Println(fmt.Sprintf("[+] Found Runtime class: id=%v", runtimeClas.ReferenceTypeID))
	methods, _ := debuggerCore.VMCommands().AllMethods(runtimeClas.ReferenceTypeID) // 10d9
	getRuntimeMethod := getMethodByName(methods, "getRuntime")
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
		if threadStatus.ThreadStatus == 2 { // thread sleeping
			threadID = thread.ObjectID
			break
		}
	}
	fmt.Println(fmt.Sprintf("[+] Setting 'step into' event in thread: %v", threadID))
	debuggerCore.VMCommands().Suspend()
	reply, _ := debuggerCore.VMCommands().SendEventRequest(1, threadID)
	fmt.Println(reply.RequestID)
	debuggerCore.VMCommands().Resume()

	//res, _ := s.ReadPacket()
	//fmt.Println(res.String())

	buf := make([]byte, 128)
	var rId int32
	var tId uint64
	num, _ := conn.Read(buf) // 好像和 read那个 会有一个阻塞 然后就会导致读不到
	if num != 0 {
		replyData := buf[:num]
		rId, tId = parseEvent(replyData, reply.RequestID, idSizes)
		fmt.Println(hex.EncodeToString(replyData))
	}
	fmt.Println(fmt.Sprintf("[+] Received matching event from thread %v", tId))
	fmt.Println(rId)
	debuggerCore.VMCommands().ClearCommand(rId)

	// 获取之后进行命令执行
	runtimeExec()

}
