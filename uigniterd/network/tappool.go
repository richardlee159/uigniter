package network

import (
	"errors"
	"log"

	"github.com/satori/uuid"
	"github.com/vishvananda/netlink"
)

const (
	DefaultBridgeName    = "uigniter0"
	DefaultBridgeAddress = "172.16.0.1/24"
	TapNamePrefix        = "utap"
	TapNumberIncrement   = 5
)

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

func createTap() (*netlink.Tuntap, error) {
	attr := netlink.NewLinkAttrs()
	attr.Name = TapNamePrefix + uuid.NewV4().String()[:8]
	tap := &netlink.Tuntap{LinkAttrs: attr, Mode: netlink.TUNTAP_MODE_TAP}

	err := netlink.LinkAdd(tap)
	return tap, err
}

type TapPool struct {
	bridge *netlink.Bridge
	taps   []*netlink.Tuntap
}

func NewTapPool() *TapPool {
	newBridge, err := getBridge()
	if err != nil {
		log.Fatalln(err)
	}
	return &TapPool{
		bridge: newBridge,
		taps:   []*netlink.Tuntap{},
	}
}

func (p *TapPool) Alloc() (string, error) {
	if len(p.taps) == 0 {
		if p.addTaps() == 0 {
			return "", errors.New("No tap device to use")
		}
	}
	tapsLen := len(p.taps)
	tap := p.taps[tapsLen-1]
	err := netlink.LinkSetUp(tap)
	if err != nil {
		return "", err
	}
	p.taps = p.taps[:tapsLen-1]
	return tap.Name, nil
}

func (p *TapPool) Release(tapName string) error {
	link, err := netlink.LinkByName(tapName)
	tap := link.(*netlink.Tuntap)
	if err != nil {
		return err
	}
	err = netlink.LinkSetDown(tap)
	if err != nil {
		return err
	}
	p.taps = append(p.taps, tap)
	return nil
}

func (p *TapPool) addTaps() int {
	count := 0
	for count < TapNumberIncrement {
		newTap, err := createTap()
		if err != nil {
			log.Println(err)
			break
		}
		err = netlink.LinkSetMaster(newTap, p.bridge)
		if err != nil {
			log.Println(err)
			break
		}
		p.taps = append(p.taps, newTap)
		count++
	}
	return count
}

func (p *TapPool) deleteAllTaps() {
	for _, tap := range p.taps {
		netlink.LinkDel(tap)
	}
	p.taps = p.taps[0:0]
}
