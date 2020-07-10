package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/richardlee159/uigniter/uigniterd/network"
)

const (
	FirecrackerBinary = "firecracker"
	RepositoryRoot    = "/var/lib/uigniter"
	ImageRoot         = RepositoryRoot + "/image"
	KernelRoot        = RepositoryRoot + "/kernel"
	VMRoot            = RepositoryRoot + "/vm"
	DefaultKernel     = KernelRoot + "/kernel.elf"

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
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	handleExit()
}

func handleExit() {
}
