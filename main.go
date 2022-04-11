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
func parseEvent(buf []byte) {

}

func main() {

	/* hacky test playground for now */

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
	}
	fmt.Printf("version = %v\n", version)

	allClasses, err := debuggercore.VMCommands().AllClasses()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	//fmt.Printf("allclasses = %v\n", allClasses)

	idSizes, err := debuggercore.VMCommands().IDSizes()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	fmt.Printf("idSizes = %v\n", idSizes)

	/**
	找到 runtime 的那个类
	*/
	for _, clas := range allClasses.Classes {
		if clas.Signature.String() == "Ljava/lang/Runtime;" {
			methods, _ := debuggercore.VMCommands().AllMethods(clas.ReferenceTypeID) // 10d9
			getRuntimeMethod := getMethodByName(methods, "getRuntime")
			if getRuntimeMethod == nil {
				return
			}
			fmt.Println(getRuntimeMethod.String())
			threads, err := debuggercore.VMCommands().AllThreads()
			if err != nil {
				fmt.Printf("err = %v\n", err)
			}
			var threadID uint64
			for _, thread := range threads.Threads {
				threadStatus, _ := debuggercore.VMCommands().StatusThread(thread.ObjectID)
				// 遍历线程的过程中需要编写一个函数去查询线程当前的状态
				if threadStatus.ThreadStatus == 2 {
					threadID = thread.ObjectID
					break
				}
				//if thread.ObjectID == uint64(4417) {
				//	fmt.Println(thread.ObjectID)
				//	threadID = thread.ObjectID
				//}
			}

			fmt.Println(threadID)

			fmt.Println(fmt.Sprintf("[+] Setting 'step into' event in thread: %v", threadID))
			//// 这里应该是下断点
			debuggercore.VMCommands().Suspend()
			reply, _ := debuggercore.VMCommands().SendEventRequest(1, threadID)
			fmt.Println(reply.RequestID)
			debuggercore.VMCommands().Resume()

			buf := make([]byte, 1024)
			for {
				num, _ := conn.Read(buf)
				fmt.Println(num)
				// 获取到了返回的 event 之后要进行解析
				fmt.Println(hex.EncodeToString(buf[:num]))
				//parseEvent(buf)
				if num != 0 {
					break
				}
			}
			debuggercore.VMCommands().ClearCommand(reply.RequestID)
		}
	}

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
