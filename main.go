package main

import (
	"encoding/hex"
	"fmt"
	"github.com/kpli0rn/jdwpgo/debuggercore"
	"github.com/kpli0rn/jdwpgo/jdwpsession"
	"github.com/kpli0rn/jdwpgo/protocol/vm"
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
func parseEvent(buf []byte) {

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
	debuggercore := debuggercore.NewFromJWDPSession(s)
	version, err := debuggercore.VMCommands().Version()
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}
	fmt.Printf("[+] version = %v\n", version)
	allClasses, err := debuggercore.VMCommands().AllClasses()
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	idSizes, err := debuggercore.VMCommands().IDSizes() // 这里的 IDsize 没有用起来
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
	methods, _ := debuggercore.VMCommands().AllMethods(runtimeClas.ReferenceTypeID) // 10d9
	getRuntimeMethod := getMethodByName(methods, "getRuntime")
	if getRuntimeMethod == nil {
		return
	}
	fmt.Println(fmt.Sprintf("[+] find getRuntime Method %s", getRuntimeMethod.String()))
	threads, err := debuggercore.VMCommands().AllThreads()
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	var threadID uint64
	for _, thread := range threads.Threads {
		threadStatus, _ := debuggercore.VMCommands().StatusThread(thread.ObjectID)
		if threadStatus.ThreadStatus == 2 { // thread sleeping
			threadID = thread.ObjectID
			break
		}
	}
	fmt.Println(fmt.Sprintf("[+] Setting 'step into' event in thread: %v", threadID))
	debuggercore.VMCommands().Suspend()
	reply, _ := debuggercore.VMCommands().SendEventRequest(1, threadID)
	debuggercore.VMCommands().Resume()
	fmt.Println(reply.RequestID)

	//var res *jdwpsession.WrappedPacket
	//for {
	//	res, _ = s.ReadPacket() // 问题主要出在这里，为什么没有返回难道是因为协程的问题嘛
	//	if res != nil {
	//		fmt.Println(res)
	//		break
	//	}
	//}
	buf := make([]byte, 128)
	for {
		num, _ := conn.Read(buf) // 好像和 read那个 会有一个阻塞 然后就会导致读不到
		if num != 0 {
			replyData := buf[:num]
			parseEvent(replyData)
			fmt.Println(hex.EncodeToString(replyData))
			break
		}
	}
	//debuggercore.VMCommands().ClearCommand(reply.RequestID)

	//caps, err := debuggercore.VMCommands().Capabilities()
	//
	//if err != nil {
	//	fmt.Printf("err = %v\n", err)
	//}
	//fmt.Printf("caps = %v\n", caps)
	//
	//capsNew, err := debuggercore.VMCommands().CapabilitiesNew()
	//
	//if err != nil {
	//	fmt.Printf("err = %v\n", err)
	//}
	//fmt.Printf("caps = %v\n", capsNew)
	//
	//tlg, err := debuggercore.VMCommands().TopLevelThreadGroups()
	//
	//if err != nil {
	//	fmt.Printf("err = %v\n", err)
	//}
	//fmt.Printf("tlgs = %v\n", tlg)
	//
	//err = debuggercore.VMCommands().Resume()
	//
	//if err != nil {
	//	fmt.Printf("err = %v\n", err)
	//}
	//
	//time.Sleep(time.Second)
	//
	//allThreads, err := debuggercore.VMCommands().AllThreads()
	//
	//if err != nil {
	//	fmt.Printf("err = %v\n", err)
	//}
	//fmt.Printf("allThreads = %v\n", allThreads)
	//
	//for idx, threadID := range allThreads.Threads {
	//	name, err := debuggercore.ThreadCommands().Name(threadID)
	//
	//	if err != nil {
	//		fmt.Printf("err = %v\n", err)
	//	}
	//	fmt.Printf("thread idx = %v, tid= %s name=%s\n", idx, threadID.String(), (string)(name.ByteString))
	//
	//}

}
