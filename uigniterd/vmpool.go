package main

import (
	"log"

	"github.com/richardlee159/uigniter/uigniterd/network"
)

const vmPoolCapacity = 1

var (
	readyVMs      chan *FirecrackerVM
	terminatedVMs chan *FirecrackerVM
	runningVMs    map[*FirecrackerVM]bool
)

func InitVMPool() {
	readyVMs = make(chan *FirecrackerVM, vmPoolCapacity)
	terminatedVMs = make(chan *FirecrackerVM, vmPoolCapacity)
	runningVMs = make(map[*FirecrackerVM]bool)
	go handleVMPool()
}

func DestroyVMPool() {
	for vm := range runningVMs {
		vm.Stop()
	}
}

func AllocVM() *FirecrackerVM {
	vm := <-readyVMs
	return vm
}

// func ReleaseVM(vm *FirecrackerVM) {
// 	terminatedVMs <- vm
// }

func createReadyVM() (*FirecrackerVM, error) {
	vm := NewVM()
	runningVMs[vm] = true

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
	delete(runningVMs, vm)
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
