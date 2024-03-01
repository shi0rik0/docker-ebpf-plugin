package main

import (
	"fmt"
	"log"

	"github.com/docker/go-plugins-helpers/network"
	"github.com/vishvananda/netlink"
)

const (
	PLUGIN_NAME = "ebpf"
)

type Driver struct {
}

func (d *Driver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	log.Printf("GetCapabilities()\n")
	return &network.CapabilitiesResponse{Scope: "local", ConnectivityScope: "local"}, nil
}

func (d *Driver) CreateNetwork(request *network.CreateNetworkRequest) error {
	log.Printf("CreateNetwork(): NetworkID: %s, Gateway: %s\n", request.NetworkID[:6], request.IPv4Data[0].Gateway)
	return nil // TODO: impl
}

func (d *Driver) AllocateNetwork(request *network.AllocateNetworkRequest) (*network.AllocateNetworkResponse, error) {
	log.Printf("AllocateNetwork(): %v\n", request)
	return nil, nil // TODO: impl
}

func (d *Driver) DeleteNetwork(request *network.DeleteNetworkRequest) error {
	log.Printf("DeleteNetwork(): %v\n", request)
	return nil // TODO: impl
}

func (d *Driver) FreeNetwork(request *network.FreeNetworkRequest) error {
	log.Printf("FreeNetwork(): %v\n", request)
	return nil // TODO: impl
}

func (d *Driver) CreateEndpoint(request *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	log.Printf("CreateEndpoint(): %v\n", request)
	return nil, nil // TODO: impl
}

func (d *Driver) DeleteEndpoint(request *network.DeleteEndpointRequest) error {
	log.Printf("DeleteEndpoint(): %v\n", request)
	return nil // TODO: impl
}

func (d *Driver) EndpointInfo(request *network.InfoRequest) (*network.InfoResponse, error) {
	log.Printf("EndpointInfo(): %v\n", request)
	return nil, nil // TODO: impl
}

func (d *Driver) Join(request *network.JoinRequest) (*network.JoinResponse, error) {
	log.Printf("Join(): %v\n", request)

	name1 := "veth-test-7"
	name2 := "veth-test-8"
	linkAttr := netlink.NewLinkAttrs()
	linkAttr.Name = name1
	veth := &netlink.Veth{LinkAttrs: linkAttr, PeerName: name2}
	err := netlink.LinkAdd(veth)
	if err != nil {
		fmt.Print(err)
	}
	if1, _ := netlink.LinkByName(name1)
	if2, _ := netlink.LinkByName(name2)
	netlink.LinkSetUp(if1)
	netlink.LinkSetUp(if2)

	return &network.JoinResponse{InterfaceName: network.InterfaceName{
		SrcName:   name2,
		DstPrefix: "eth",
	}}, nil // TODO: impl
}

func (d *Driver) Leave(request *network.LeaveRequest) error {
	// 只需要删除host一侧的veth网卡即可。
	log.Printf("Leave(): %v\n", request)
	return nil // TODO: impl
}

func (d *Driver) DiscoverNew(request *network.DiscoveryNotification) error {
	log.Printf("DiscoverNew(): %v\n", request)
	return nil // TODO: impl
}

func (d *Driver) DiscoverDelete(request *network.DiscoveryNotification) error {
	log.Printf("DiscoverDelete(): %v\n", request)
	return nil // TODO: impl
}

func (d *Driver) ProgramExternalConnectivity(request *network.ProgramExternalConnectivityRequest) error {
	log.Printf("ProgramExternalConnectivity(): %v\n", request)
	return nil // TODO: impl
}

func (d *Driver) RevokeExternalConnectivity(request *network.RevokeExternalConnectivityRequest) error {
	log.Printf("RevokeExternalConnectivity(): %v\n", request)
	return nil // TODO: impl
}

func main() {

	handler := network.NewHandler(&Driver{})
	err := handler.ServeUnix(PLUGIN_NAME, 0)
	if err != nil {
		log.Fatal(err)
	}
}
