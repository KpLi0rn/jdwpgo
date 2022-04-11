package main

import (
	"fmt"
	"github.com/jquirke/jdwpgo/debuggercore"
	"github.com/jquirke/jdwpgo/jdwpsession"
	"github.com/jquirke/jdwpgo/protocol/vm"
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
			//fmt.Println(clas.ReferenceTypeID.String())
			methods, _ := debuggercore.VMCommands().AllMethods(clas.ReferenceTypeID)
			//fmt.Println(methods.String())
			getRuntimeMethod := getMethodByName(methods, "getRuntime")
			if getRuntimeMethod == nil {
				return
			}

			fmt.Println(getRuntimeMethod.String())

			threads, err := debuggercore.VMCommands().AllThreads()
			if err != nil {
				fmt.Printf("err = %v\n", err)
			}

			var threadId uint64
			for _, thread := range threads.Threads {
				// 这里先固定一下看一下后面走不走的通
				if thread.ObjectID == uint64(4471) {
					threadId = thread.ObjectID
				}
			}
			fmt.Println(fmt.Sprintf("[+] Setting 'step into' event in thread: %v", threadId))

			// 这里应该是下断点
			debuggercore.VMCommands().Suspend()
			debuggercore.VMCommands().Resume()
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
