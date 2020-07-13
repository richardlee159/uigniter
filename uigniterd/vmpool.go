package main

import (
	"errors"
	"log"

	"github.com/richardlee159/uigniter/uigniterd/network"
)

const vmPoolCapacity = 1

var (
	readyVMs      chan *FirecrackerVM
	terminatedVMs chan *FirecrackerVM
	runningVMs    map[string]*FirecrackerVM
)

func InitVMPool() {
	readyVMs = make(chan *FirecrackerVM, vmPoolCapacity)
	terminatedVMs = make(chan *FirecrackerVM, vmPoolCapacity)
	runningVMs = make(map[string]*FirecrackerVM)
	go handleVMPool()
}

func DestroyVMPool() {
	for _, vm := range runningVMs {
		vm.Stop()
	}
}

func RunVM(opt *Options) (*FirecrackerVM, error) {
	vm := <-readyVMs

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
	vm, ok := runningVMs[id]
	if !ok {
		return errors.New("VM not found")
	}
	return vm.Stop()
}

func createReadyVM() (*FirecrackerVM, error) {
	vm := NewVM()
	runningVMs[vm.uuid] = vm

	go func() {
		err := vm.Wait()
		if err != nil {
			log.Print(err)
		}
		terminatedVMs <- vm
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

	vm.tapName = tapName
	vm.ipAddr = ipAddr
	err = vm.ConfigNetwork(tapName, macAddr)

	return vm, err
}

func deleteTermVM(vm *FirecrackerVM) {
	delete(runningVMs, vm.uuid)
	err := network.DeleteTap(vm.tapName)
	if err != nil {
		log.Println(err)
	}
	network.ReleaseIPv4(vm.ipAddr)
	vm.Delete()
}

func handleVMPool() {
	rvm, _ := createReadyVM()
	for {
		select {
		case readyVMs <- rvm:
			rvm, _ = createReadyVM()
		case tvm := <-terminatedVMs:
			deleteTermVM(tvm)
		}
	}
}
