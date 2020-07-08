package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type FirecrackerVM struct {
	name string
	uuid string
	conf Config
}

func (vm *FirecrackerVM) ConfigBasic(opt Options) {
	vm.conf.MachineCfg = MachineConfig{
		VcpuCount: opt.VcpuCount,
		MemSize:   opt.MemSize,
		HtEnabled: false,
	}

	vm.conf.BootSource = BootSource{
		KernelPath: opt.KernelPath,
		BootArgs:   "--nopci " + opt.CommandLine,
	}

	vm.conf.Drives = []Drive{
		Drive{
			DriveId:    "rootfs",
			DiskPath:   opt.DiskPath,
			RootDevice: false,
			ReadOnly:   opt.ReadOnly,
		},
	}
}

func (vm *FirecrackerVM) ConfigNetwork() {

}

func (vm *FirecrackerVM) Start() {
	tmpfile, err := ioutil.TempFile("", "uigniter-")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(vm.conf.GetJson()); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(FirecrackerBinary, "--no-api", "--config-file", tmpfile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
