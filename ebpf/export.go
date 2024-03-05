package ebpf

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/shi0rik0/docker-ebpf-plugin/tc"
	"github.com/vishvananda/netlink"
)

type Program struct {
	tcObjects
}

func LoadProgram() (*Program, error) {
	var objs tcObjects
	err := loadTcObjects(&objs, nil)
	if err != nil {
		return nil, err
	}
	program := &Program{tcObjects: objs}
	return program, nil
}

func (p *Program) Close() error {
	return p.tcObjects.Close()
}

func (p *Program) Attach(interfaceName string, programName string) error {
	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return err
	}
	return tc.AttachTC(p.TcIngress, link, tc.INGRESS, programName)
}

func (p *Program) Detach(interfaceName string) error {
	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return err
	}
	return tc.DetachTC(p.TcIngress, link, tc.INGRESS)
}

func (p *Program) AddIpIfindexMapEntry(ipAddr string, ifindex uint32) error {
	ip, err := ipAddrToUint32(ipAddr)
	if err != nil {
		return err
	}
	return p.tcMaps.IpIfindexMap.Put(ip, ifindex)
}

func (p *Program) DeleteIpIfindexMapEntry(ipAddr string) error {
	ip, err := ipAddrToUint32(ipAddr)
	if err != nil {
		return err
	}
	return p.tcMaps.IpIfindexMap.Delete(ip)
}

func ipAddrToUint32(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("invalid IP address")
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0, errors.New("not an IPv4 address")
	}

	return binary.LittleEndian.Uint32(ipv4), nil
}
