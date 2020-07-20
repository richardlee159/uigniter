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
	"syscall"
	"time"

	"github.com/satori/uuid"
)

type Firecracker struct {
	uuid      string
	tapName   string
	ipAddr    string
	conf      *Config
	cmd       *exec.Cmd
	apiClient *http.Client
}

func NewFC() *Firecracker {
	fc := &Firecracker{
		uuid: uuid.NewV4().String(),
		conf: &Config{},
	}

	fc.cmd = exec.Command(FirecrackerBinary, "--id", fc.uuid, "--api-sock", fc.socketPath())
	fc.cmd.Stdin = os.Stdin
	fc.cmd.Stdout = os.Stdout
	fc.cmd.Stderr = os.Stderr
	err := fc.cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(fc.socketPath())
	for os.IsNotExist(err) {
		time.Sleep(100 * time.Millisecond)
		_, err = os.Stat(fc.socketPath())
	}

	fc.apiClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", fc.socketPath())
			},
		},
	}
	return fc
}

func (fc *Firecracker) Delete() {
	os.Remove(fc.socketPath())
}

func (fc *Firecracker) socketPath() string {
	return VMRoot + fc.uuid + ".socket"
}

func (fc *Firecracker) apiPutCall(path string, data interface{}) error {
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

	response, err := fc.apiClient.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New(response.Status)
	}

	return nil
}

func (fc *Firecracker) ConfigMachine(vcpus, memSize int) error {
	fc.conf.MachineCfg = MachineConfig{
		VcpuCount: vcpus,
		MemSize:   memSize,
		HtEnabled: false,
	}
	return fc.apiPutCall("/machine-config", fc.conf.MachineCfg)
}

func (fc *Firecracker) ConfigBootSource(kernel, bootArgs string) error {
	fc.conf.BootSource = BootSource{
		KernelPath: kernel,
		BootArgs:   bootArgs,
	}
	return fc.apiPutCall("/boot-source", fc.conf.BootSource)
}

func (fc *Firecracker) ConfigRootfs(disk string, readonly bool) error {
	fc.conf.Drives = []Drive{
		Drive{
			DriveId:    "rootfs",
			DiskPath:   disk,
			RootDevice: false,
			ReadOnly:   readonly,
		},
	}
	return fc.apiPutCall("/drives/rootfs", fc.conf.Drives[0])
}

func (fc *Firecracker) ConfigNetwork(tapName, macAddr string) error {
	fc.conf.NwInterfaces = []NetworkInterface{
		NetworkInterface{
			IfaceId: "eth0",
			TapName: tapName,
			Mac:     macAddr,
		},
	}
	return fc.apiPutCall("/network-interfaces/eth0", fc.conf.NwInterfaces[0])
}

func (fc *Firecracker) Start() error {
	data := make(map[string]string)
	data["action_type"] = "InstanceStart"
	return fc.apiPutCall("/actions", data)
}

func (vm *Firecracker) Wait() error {
	return vm.cmd.Wait()
}

func (vm *Firecracker) Stop() error {
	return vm.cmd.Process.Signal(syscall.SIGINT)
}
