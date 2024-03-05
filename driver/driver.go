package main

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/docker/go-plugins-helpers/network"
	"github.com/shi0rik0/docker-ebpf-plugin/ebpf"
	"github.com/vishvananda/netlink"
)

const (
	PLUGIN_NAME      = "ebpf"
	HOST_VETH_PREFIX = "veth-ebpf-"
)

type Driver struct {
	mutex    sync.Mutex
	vethName map[connectionID]string
	program  map[string]*ebpf.Program
	counter  int
}

func NewDriver() *Driver {
	return &Driver{
		vethName: make(map[connectionID]string),
		program:  make(map[string]*ebpf.Program),
	}
}

type connectionID struct {
	NetworkID  string
	EndpointID string
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

	d.mutex.Lock()
	defer d.mutex.Unlock()

	cid := connectionID{NetworkID: request.NetworkID, EndpointID: request.EndpointID}
	if _, ok := d.vethName[cid]; ok {
		return nil, errors.New("connection exists")
	}

	id := strconv.Itoa(d.counter)
	name1 := HOST_VETH_PREFIX + id
	name2 := "veth-temp-" + id
	err := createVeth(name1, name2)
	if err != nil {
		return nil, err
	}

	d.vethName[cid] = name1
	d.counter++
	return &network.JoinResponse{InterfaceName: network.InterfaceName{
		SrcName:   name2,
		DstPrefix: "eth",
	}}, nil // TODO: impl
}

func createVeth(name1 string, name2 string) error {
	linkAttr := netlink.NewLinkAttrs()
	linkAttr.Name = name1
	veth := &netlink.Veth{LinkAttrs: linkAttr, PeerName: name2}
	err := netlink.LinkAdd(veth)
	if err != nil {
		return err
	}
	if1, err := netlink.LinkByName(name1)
	if err != nil {
		return err
	}
	if2, err := netlink.LinkByName(name2)
	if err != nil {
		return err
	}
	netlink.LinkSetUp(if1)
	netlink.LinkSetUp(if2)
	return nil
}

func deleteVeth(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	err = netlink.LinkDel(link)
	if err != nil {
		return err
	}
	return nil
}

func (d *Driver) Leave(request *network.LeaveRequest) error {
	log.Printf("Leave(): %v\n", request)

	d.mutex.Lock()
	defer d.mutex.Unlock()

	cid := connectionID{NetworkID: request.NetworkID, EndpointID: request.EndpointID}
	name, ok := d.vethName[cid]
	if !ok {
		return errors.New("connection doesn't exist")
	}

	err := deleteVeth(name)
	if err != nil {
		return err
	}

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
	driver := NewDriver()
	handler := network.NewHandler(driver)
	err := handler.ServeUnix(PLUGIN_NAME, 0)
	if err != nil {
		log.Fatal(err)
	}
}
