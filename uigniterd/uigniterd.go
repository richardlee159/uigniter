package main

import (
	"os"
	"os/signal"
	"syscall"

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

}

func main() {
	exitsig := make(chan os.Signal)
	signal.Notify(exitsig, syscall.SIGINT, syscall.SIGTERM)

	childterm := make(chan os.Signal)
	signal.Notify(childterm, syscall.SIGCHLD)

	InitVMPool()
	runtest()

	<-exitsig
	handleExit()
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
}
