package main

import (
	"log"

	"github.com/richardlee159/uigniter/uigniterd/network"
)

const vmPoolCapacity = 1

var (
	readyVMs      chan *FirecrackerVM
	terminatedVMs chan *FirecrackerVM
)

func InitVMPool() {
	readyVMs = make(chan *FirecrackerVM, vmPoolCapacity)
	terminatedVMs = make(chan *FirecrackerVM)
	go handleVMPool()
}

func DestroyVMPool() {

}

func AllocVM() (*FirecrackerVM, error) {
	vm := <-readyVMs
	return vm, nil
}

func ReleaseVM(vm *FirecrackerVM) {
	terminatedVMs <- vm
}

func createReadyVM() (*FirecrackerVM, error) {
	vm := NewVM()

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
	err := network.DeleteTap(vm.tapName)
	if err != nil {
		log.Println(err)
	}
	network.ReleaseIPv4(vm.ipAddr)
	vm.Delete()
}

func handleVMPool() {
	var vm *FirecrackerVM
	vm, _ = createReadyVM()
	for {
		select {
		case readyVMs <- vm:
			vm, _ = createReadyVM()
		case vm := <-terminatedVMs:
			deleteTermVM(vm)
		}
	}
}
