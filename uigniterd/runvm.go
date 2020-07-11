package main

import "log"

func runVM(opt *Options) error {
	vm, err := AllocVM()
	if err != nil {
		return err
	}

	go func() {
		err := vm.cmd.Wait()
		if err != nil {
			log.Print(err)
		}
		log.Println("release:", vm.uuid)
		ReleaseVM(vm)
	}()

	bootArgs := "--nopci" +
		" --ip=eth0," + vm.ipAddr + "," + DefaultSubnetMask +
		" --defaultgw=" + DefaultGateway +
		" --nameserver=" + DefaultNameServer + " " +
		opt.CommandLine
	err = vm.ConfigBootSource(opt.KernelPath, bootArgs)
	if err != nil {
		return err
	}
	err = vm.ConfigRootfs(opt.DiskPath, opt.ReadOnly)
	if err != nil {
		return err
	}

	return vm.Start()
}
