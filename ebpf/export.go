package ebpf

import (
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
