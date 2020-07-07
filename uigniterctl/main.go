package main

const (
	FirecrackerBinary = "firecracker"
	RepositoryRoot    = "/var/lib/uigniter/"
	DefaultKernel     = "kernel/kernel.elf"
	DefaultGateway    = "172.16.0.0"
	DefaultSubnetMask = "255.255.0.0"
)

type Options struct {
	VcpuCount   int
	MemSize     int
	KernelPath  string
	DiskPath    string
	ReadOnly    bool
	CommandLine string
}

func main() {
	opt := Options{
		1,
		128,
		RepositoryRoot + DefaultKernel,
		RepositoryRoot + "vm/hello/hello.raw",
		true,
		"--bootchart hello",
	}
	vm0 := FirecrackerVM{}
	vm0.ConfigBasic(opt)
	vm0.Start()
}
