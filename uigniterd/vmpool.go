package main

import (
	"errors"
	"log"

	"github.com/richardlee159/uigniter/uigniterd/network"
)

const vmPoolCapacity = 1

var (
	readyFCs      chan *Firecracker
	terminatedFCs chan *Firecracker
	runningFCs    map[string]*Firecracker
)

func InitVMPool() {
	readyFCs = make(chan *Firecracker, vmPoolCapacity)
	terminatedFCs = make(chan *Firecracker, vmPoolCapacity)
	runningFCs = make(map[string]*Firecracker)
	go handleVMPool()
}

func DestroyVMPool() {
	for _, vm := range runningFCs {
		vm.Stop()
	}
}

func RunVM(opt *Options) (*Firecracker, error) {
	vm := <-readyFCs

	bootArgs := "--nopci" +
		" --ip=eth0," + vm.ipAddr + "," + DefaultSubnetMask +
		" --defaultgw=" + DefaultGateway +
		" --nameserver=" + DefaultNameServer + " " +
		opt.CommandLine
	err := vm.ConfigBootSource(opt.KernelPath, bootArgs)
	if err != nil {
		vm.Stop()
		return nil, err
	}
	err = vm.ConfigRootfs(opt.DiskPath, opt.ReadOnly)
	if err != nil {
		vm.Stop()
		return nil, err
	}
	err = vm.Start()
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func StopVM(id string) error {
	vm, ok := runningFCs[id]
	if !ok {
		return errors.New("VM not found")
	}
	return vm.Stop()
}

func createReadyFC() (*Firecracker, error) {
	fc := NewFC()
	runningFCs[fc.uuid] = fc

	go func() {
		err := fc.Wait()
		if err != nil {
			log.Print(err)
		}
		terminatedFCs <- fc
	}()

	tapName, err := network.NewTap()
	if err != nil {
		return nil, err
	}
	ipAddr, err := network.AllocIPv4()
	if err != nil {
		return nil, err
	}
	macAddr, err := network.GenMac()
	if err != nil {
		return nil, err
	}

	fc.tapName = tapName
	fc.ipAddr = ipAddr
	err = fc.ConfigNetwork(tapName, macAddr)

	return fc, err
}

func deleteTermVM(vm *Firecracker) {
	delete(runningFCs, vm.uuid)
	err := network.DeleteTap(vm.tapName)
	if err != nil {
		log.Println(err)
	}
	network.ReleaseIPv4(vm.ipAddr)
	vm.Delete()
}

func handleVMPool() {
	rvm, _ := createReadyFC()
	for {
		select {
		case readyFCs <- rvm:
			rvm, _ = createReadyFC()
		case tvm := <-terminatedFCs:
			deleteTermVM(tvm)
		}
	}
}
