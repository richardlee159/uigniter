package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/richardlee159/uigniter/uigniterd/network"
)

const (
	FirecrackerBinary = "firecracker"
	RepositoryRoot    = "/var/lib/uigniter/"
	ImageRoot         = RepositoryRoot + "image/"
	KernelRoot        = RepositoryRoot + "kernel/"
	VMRoot            = RepositoryRoot + "vm/"
	DefaultKernel     = KernelRoot + "kernel.elf"

	DefaultGateway    = network.DefaultBridgeIP
	DefaultNameServer = network.DefaultBridgeIP
	DefaultSubnetMask = "255.255.225.0"
)

type Options struct {
	VcpuCount   int
	MemSize     int
	KernelPath  string
	DiskPath    string
	ReadOnly    bool
	CommandLine string
}

func init() {
	checkRepo()
	InitVMPool()
}

func main() {
	exitsig := make(chan os.Signal)
	signal.Notify(exitsig, syscall.SIGINT, syscall.SIGTERM)

	childterm := make(chan os.Signal)
	signal.Notify(childterm, syscall.SIGCHLD)

	runtest()

	<-exitsig
	handleExit()
}

func checkRepo() {
	_, err := os.Stat(RepositoryRoot)
	if os.IsNotExist(err) {
		fmt.Println("Initializing repository...")
		err = os.MkdirAll(RepositoryRoot, 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Mkdir(KernelRoot, 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Mkdir(ImageRoot, 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Mkdir(VMRoot, 0755)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Need to add kernel and image manually!")
	} else if err != nil {
		log.Print(err)
	}
}

// func handleChild() {
// 	var wstatus syscall.WaitStatus
// 	for {
// 		pid, err := syscall.Wait4(-1, &wstatus, syscall.WNOHANG, nil)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}

// }

func handleExit() {
	DestroyVMPool()
	time.Sleep(time.Second)
}
