package main

func runVM(opt *Options) error {
	vm, err := AllocVM()
	if err != nil {
		return err
	}

	bootArgs := "--ip=eth0," + vm.ipAddr + "," + DefaultSubnetMask +
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
