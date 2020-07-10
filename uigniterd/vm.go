package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/satori/uuid"
)

type FirecrackerVM struct {
	uuid    string
	conf    *Config
	cmd     *exec.Cmd
	apiCli  *http.Client
	tapName string
	ipAddr  string
}

func NewVM() *FirecrackerVM {
	vm := &FirecrackerVM{
		uuid: uuid.NewV4().String(),
		conf: &Config{},
	}

	vm.cmd = exec.Command(FirecrackerBinary, "--id", vm.uuid, "--api-sock", vm.socketPath())
	vm.cmd.Stdout = os.Stdout
	vm.cmd.Stderr = os.Stderr
	err := vm.cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	vm.apiCli = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", vm.socketPath())
			},
		},
	}
	return vm
}

func (vm *FirecrackerVM) socketPath() string {
	return VMRoot + vm.uuid + ".socket"
}

func (vm *FirecrackerVM) apiPutCall(path string, data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut,
		"http://localhost"+path, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	response, err := vm.apiCli.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New(response.Status)
	}

	return nil
}

func (vm *FirecrackerVM) ConfigMachine(vcpus, memSize int) error {
	vm.conf.MachineCfg = MachineConfig{
		VcpuCount: vcpus,
		MemSize:   memSize,
		HtEnabled: false,
	}
	return vm.apiPutCall("/machine-config", vm.conf.MachineCfg)
}

func (vm *FirecrackerVM) ConfigBootSource(kernel, cmdline string) error {
	vm.conf.BootSource = BootSource{
		KernelPath: kernel,
		BootArgs:   cmdline,
	}
	return vm.apiPutCall("/boot-source", vm.conf.MachineCfg)
}

func (vm *FirecrackerVM) ConfigRootfs(disk string, readonly bool) error {
	vm.conf.Drives = []Drive{
		Drive{
			DriveId:    "rootfs",
			DiskPath:   disk,
			RootDevice: true,
			ReadOnly:   readonly,
		},
	}
	return vm.apiPutCall("/drives/rootfs", vm.conf.Drives[0])
}

func (vm *FirecrackerVM) ConfigNetwork(tapName, macAddr string) error {
	vm.conf.NwInterfaces = []NetworkInterface{
		NetworkInterface{
			IfaceId: "eth0",
			TapName: tapName,
			Mac:     macAddr,
		},
	}
	return vm.apiPutCall("network-interfaces/eth0", vm.conf.NwInterfaces[0])
}

func (vm *FirecrackerVM) Start() error {
	var data map[string]string
	data["action_type"] = "InstanceStart"
	return vm.apiPutCall("/actions", data)
}