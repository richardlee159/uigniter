package main

import (
	"errors"
	"log"

	"github.com/richardlee159/uigniter/uigniterd/network"
)

const vmPoolCapacity = 5

var vmPool []*FirecrackerVM

func AllocVM() (*FirecrackerVM, error) {
	l := len(vmPool)
	if l == 0 {
		err := createReadyVMs()
		l = len(vmPool)
		if l == 0 {
			return nil, errors.New("No vm to alloc")
		}
		if err != nil {
			log.Println(err)
		}
	}
	vm := vmPool[l-1]
	vmPool = vmPool[:l-1]
	return vm, nil
}

func ReleaseVM(vm *FirecrackerVM) {
	err := network.DeleteTap(vm.tapName)
	if err != nil {
		log.Println(err)
	}
	network.ReleaseIPv4(vm.ipAddr)
}

func createReadyVMs() error {
	for i := 0; i < vmPoolCapacity; i++ {
		err := newReadyVM()
		if err != nil {
			return err
		}
	}
	return nil
}

func newReadyVM() error {
	vm := NewVM()

	tapName, err := network.NewTap()
	if err != nil {
		return err
	}
	ipAddr, err := network.AllocIPv4()
	if err != nil {
		return err
	}
	macAddr, err := network.GenMac()
	if err != nil {
		return err
	}

	vm.tapName = tapName
	vm.ipAddr = ipAddr
	err = vm.ConfigNetwork(tapName, macAddr)
	vmPool = append(vmPool, vm)

	return err
}
