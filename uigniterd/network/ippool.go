package network

import (
	"errors"
	"strconv"
	"strings"
)

const (
	DefaultNetworkPrefix = "172.17.0."
)

var hostAddrPool []int

func init() {
	hostAddrPool = newHostAddrPool()
}

func newHostAddrPool() []int {
	list := make([]int, 0, 256)
	for i := 0; i < 253; i++ {
		list = append(list, 254-i)
	}
	return list
}

// AllocIPv4 ... Alloc a new IPv4 address
func AllocIPv4() (string, error) {
	length := len(hostAddrPool)
	if length == 0 {
		return "", errors.New("No IP address to allocate")
	}
	hostAddr := hostAddrPool[length-1]
	hostAddrPool = hostAddrPool[:length-1]
	return DefaultNetworkPrefix + strconv.Itoa(hostAddr), nil
}

// ReleaseIPv4 .. Release an IPv4 address
func ReleaseIPv4(addr string) {
	hostAddr, _ := strconv.Atoi(strings.Split(addr, ".")[3])
	hostAddrPool = append(hostAddrPool, hostAddr)
}
