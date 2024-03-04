package tc

import (
	"github.com/cilium/ebpf"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type AttachPoint int

const (
	INGRESS AttachPoint = 1
	EGRESS  AttachPoint = 2
)

func AttachTC(program *ebpf.Program, link netlink.Link, attachPoint AttachPoint, name string) error {
	err := createQdisc(link.Attrs().Index)
	if err != nil {
		return err
	}

	filter := &netlink.BpfFilter{
		FilterAttrs:  makeFilterAttrs(link, attachPoint),
		Fd:           program.FD(),
		Name:         name,
		DirectAction: true,
	}
	return netlink.FilterAdd(filter)
}

func DetachTC(program *ebpf.Program, link netlink.Link, attachPoint AttachPoint) error {
	filter := &netlink.BpfFilter{
		FilterAttrs: makeFilterAttrs(link, attachPoint),
	}
	err := netlink.FilterDel(filter)
	if err != nil {
		return err
	}
	return deleteQdisc(link.Attrs().Index)
}

func makeFilterAttrs(link netlink.Link, attachPoint AttachPoint) netlink.FilterAttrs {
	var parent uint32
	switch attachPoint {
	case INGRESS:
		parent = netlink.MakeHandle(0xFFFF, 0xFFF2)
	case EGRESS:
		parent = netlink.MakeHandle(0xFFFF, 0xFFF3)
	}

	return netlink.FilterAttrs{
		LinkIndex: link.Attrs().Index,
		Handle:    1,
		Parent:    parent,
		Priority:  1,
		Protocol:  unix.ETH_P_ALL,
	}
}

func makeQdiscInfo(ifIndex int) *netlink.GenericQdisc {
	attrs := netlink.QdiscAttrs{
		LinkIndex: ifIndex,
		Handle:    netlink.MakeHandle(0xFFFF, 0),
		Parent:    0xFFFFFFF1,
	}

	return &netlink.GenericQdisc{
		QdiscAttrs: attrs,
		QdiscType:  "clsact",
	}
}

func createQdisc(ifIndex int) error {
	return netlink.QdiscAdd(makeQdiscInfo(ifIndex))
}

func deleteQdisc(ifIndex int) error {
	return netlink.QdiscDel(makeQdiscInfo(ifIndex))
}
