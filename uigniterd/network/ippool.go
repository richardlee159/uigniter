package network

import (
	"errors"
	"strconv"
	"strings"
)

const (
	DefaultNetworkPrefix = "172.16.0."
)

type IPPool struct {
	hostAddrList []int
}

func NewIPPool() *IPPool {
	list := make([]int, 0, 256)
	for i := 0; i < 253; i++ {
		list = append(list, 254-i)
	}
	return &IPPool{list}
}

func (p *IPPool) Alloc() (string, error) {
	length := len(p.hostAddrList)
	if length == 0 {
		return "", errors.New("No IP address to allocate")
	}
	hostAddr := p.hostAddrList[length-1]
	p.hostAddrList = p.hostAddrList[:length-1]
	return DefaultNetworkPrefix + strconv.Itoa(hostAddr), nil
}

func (p *IPPool) Release(addr string) {
	hostAddr, _ := strconv.Atoi(strings.Split(addr, ".")[3])
	p.hostAddrList = append(p.hostAddrList, hostAddr)
}
