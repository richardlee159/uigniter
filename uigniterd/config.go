package main

import (
	"encoding/json"
	"log"
)

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

func (conf *Config) GetJson() []byte {
	b, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
