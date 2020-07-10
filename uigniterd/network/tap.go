package network

import (
	"log"

	"github.com/satori/uuid"
	"github.com/vishvananda/netlink"
)

const (
	DefaultBridgeName    = "uigniter0"
	DefaultBridgeIP      = "172.17.0.1"
	DefaultBridgeAddress = DefaultBridgeIP + "/24"
	TapNamePrefix        = "utap"
)

var bridge *netlink.Bridge

func init() {
	var err error
	bridge, err = getBridge()
	if err != nil {
		log.Fatalln(err)
	}
}

func createBridge() (*netlink.Bridge, error) {
	attr := netlink.NewLinkAttrs()
	attr.Name = DefaultBridgeName
	bridge := &netlink.Bridge{LinkAttrs: attr}

	err := netlink.LinkAdd(bridge)
	if err != nil {
		return nil, err
	}
	addr, err := netlink.ParseAddr(DefaultBridgeAddress)
	if err != nil {
		return nil, err
	}
	err = netlink.AddrAdd(bridge, addr)
	if err != nil {
		return nil, err
	}
	err = netlink.LinkSetUp(bridge)
	if err != nil {
		return nil, err
	}
	return bridge, nil
}

func getBridge() (*netlink.Bridge, error) {
	link, err := netlink.LinkByName(DefaultBridgeName)
	if err != nil {
		return createBridge()
	}
	return link.(*netlink.Bridge), nil
}

// NewTap ... Create a new tap device and return its name
func NewTap() (string, error) {
	attr := netlink.NewLinkAttrs()
	attr.Name = TapNamePrefix + uuid.NewV4().String()[:8]
	tap := &netlink.Tuntap{LinkAttrs: attr, Mode: netlink.TUNTAP_MODE_TAP}

	err := netlink.LinkAdd(tap)
	if err != nil {
		return "", err
	}
	err = netlink.LinkSetMaster(tap, bridge)
	if err != nil {
		return "", err
	}
	err = netlink.LinkSetUp(tap)
	if err != nil {
		return "", err
	}

	return tap.Name, nil
}

// DeleteTap ... Delete a tap device by name
func DeleteTap(tapName string) error {
	link, err := netlink.LinkByName(tapName)
	if err != nil {
		return err
	}
	tap := link.(*netlink.Tuntap)

	err = netlink.LinkSetDown(tap)
	if err != nil {
		return err
	}
	err = netlink.LinkDel(tap)
	if err != nil {
		return err
	}

	return nil
}
