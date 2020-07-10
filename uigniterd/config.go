package main

type MachineConfig struct {
	VcpuCount int  `json:"vcpu_count"`
	MemSize   int  `json:"mem_size_mib"`
	HtEnabled bool `json:"ht_enabled"`
}

type BootSource struct {
	KernelPath string `json:"kernel_image_path"`
	BootArgs   string `json:"boot_args"`
}

type Drive struct {
	DriveId    string `json:"drive_id"`
	DiskPath   string `json:"path_on_host"`
	RootDevice bool   `json:"is_root_device"`
	ReadOnly   bool   `json:"is_read_only"`
}

type NetworkInterface struct {
	IfaceId string `json:"iface_id"`
	TapName string `json:"host_dev_name"`
	Mac     string `json:"guest_mac"`
}

type Config struct {
	MachineCfg   MachineConfig      `json:"machine-config"`
	BootSource   BootSource         `json:"boot-source"`
	Drives       []Drive            `json:"drives"`
	NwInterfaces []NetworkInterface `json:"network-interfaces,omitempty"`
}

func (conf *Config) AddBasicInfo(opt *Options) {
	conf.MachineCfg = MachineConfig{
		VcpuCount: opt.VcpuCount,
		MemSize:   opt.MemSize,
		HtEnabled: false,
	}

	conf.BootSource = BootSource{
		KernelPath: opt.KernelPath,
		BootArgs:   "--nopci " + opt.CommandLine,
	}

	conf.Drives = []Drive{
		Drive{
			DriveId:    "rootfs",
			DiskPath:   opt.DiskPath,
			RootDevice: false,
			ReadOnly:   opt.ReadOnly,
		},
	}
}

func (conf *Config) AddNetworkInfo(tapName, ipAddr, macAddr string) {
	conf.BootSource.BootArgs = "--ip=eth0," + ipAddr + "," + DefaultSubnetMask +
		" --defaultgw=" + DefaultGateway +
		" --nameserver=" + DefaultNameServer + " " +
		conf.BootSource.BootArgs

	conf.NwInterfaces = []NetworkInterface{
		NetworkInterface{
			IfaceId: "eth0",
			TapName: tapName,
			Mac:     macAddr,
		},
	}
}
