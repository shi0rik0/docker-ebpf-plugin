package driver

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/go-plugins-helpers/network"
	"github.com/shi0rik0/docker-ebpf-plugin/ebpf"
	"github.com/vishvananda/netlink"
)

const HOST_VETH_PREFIX = "veth-ebpf-"

type Driver struct {
	mutex             sync.Mutex
	connectionInfoMap map[connectionID]connectionInfo
	programMap        map[string]*ebpf.Program
	counter           int
}

func NewDriver() *Driver {
	return &Driver{
		connectionInfoMap: make(map[connectionID]connectionInfo),
		programMap:        make(map[string]*ebpf.Program),
	}
}

type connectionID struct {
	NetworkID  string
	EndpointID string
}

type connectionInfo struct {
	hostVethName string
	containerIp  string
}

func (d *Driver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	log.Printf("GetCapabilities()\n")

	return &network.CapabilitiesResponse{Scope: "local", ConnectivityScope: "local"}, nil
}

func (d *Driver) CreateNetwork(request *network.CreateNetworkRequest) error {
	log.Printf("CreateNetwork(): NetworkID: %s, Gateway: %s\n", request.NetworkID[:6], request.IPv4Data[0].Gateway)

	d.mutex.Lock()
	defer d.mutex.Unlock()

	p, err := ebpf.LoadProgram()
	if err != nil {
		return err
	}
	d.programMap[request.NetworkID] = p
	return nil
}

func (d *Driver) AllocateNetwork(request *network.AllocateNetworkRequest) (*network.AllocateNetworkResponse, error) {
	log.Printf("AllocateNetwork(): %v\n", request)
	return nil, nil // TODO: impl
}

func (d *Driver) DeleteNetwork(request *network.DeleteNetworkRequest) error {
	log.Printf("DeleteNetwork(): NetworkID: %s\n", request.NetworkID[:6])

	d.mutex.Lock()
	defer d.mutex.Unlock()

	p, ok := d.programMap[request.NetworkID]
	if !ok {
		return errors.New("couldn't find network")
	}
	err := p.Close()
	if err != nil {
		return err
	}
	delete(d.programMap, request.NetworkID)
	return nil
}

func (d *Driver) FreeNetwork(request *network.FreeNetworkRequest) error {
	log.Printf("FreeNetwork(): %v\n", request)
	return nil // TODO: impl
}

func (d *Driver) CreateEndpoint(request *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	log.Printf("CreateEndpoint(): NetworkID: %s, EndpointID: %s, IP: %s\n",
		request.NetworkID[:6], request.EndpointID[:6], request.Interface.Address)

	d.mutex.Lock()
	defer d.mutex.Unlock()

	cid := connectionID{NetworkID: request.NetworkID, EndpointID: request.EndpointID}
	if _, ok := d.connectionInfoMap[cid]; ok {
		return nil, errors.New("connection exists")
	}
	d.connectionInfoMap[cid] = connectionInfo{containerIp: removeSubnet(request.Interface.Address)}
	return nil, nil
}

func (d *Driver) DeleteEndpoint(request *network.DeleteEndpointRequest) error {
	log.Printf("DeleteEndpoint(): NetworkID: %s, EndpointID: %s\n", request.NetworkID[:6], request.EndpointID[:6])

	d.mutex.Lock()
	defer d.mutex.Unlock()

	cid := connectionID{NetworkID: request.NetworkID, EndpointID: request.EndpointID}
	info, ok := d.connectionInfoMap[cid]
	if !ok {
		return errors.New("connection doesn't exist")
	}

	program, ok := d.programMap[request.NetworkID]
	if !ok {
		return errors.New("couldn't find network")
	}

	err := program.DeleteIpIfindexMapEntry(info.containerIp)
	if err != nil {
		return err
	}

	err = program.Detach(info.hostVethName)
	if err != nil {
		return err
	}

	err = deleteVeth(info.hostVethName)
	if err != nil {
		return err
	}

	delete(d.connectionInfoMap, cid)
	return nil
}

func (d *Driver) EndpointInfo(request *network.InfoRequest) (*network.InfoResponse, error) {
	log.Printf("EndpointInfo(): %v\n", request)
	return nil, nil // TODO: impl
}

func (d *Driver) Join(request *network.JoinRequest) (*network.JoinResponse, error) {
	log.Printf("Join(): NetworkID: %s, EndpointID: %s\n", request.NetworkID[:6], request.EndpointID[:6])

	d.mutex.Lock()
	defer d.mutex.Unlock()

	id := strconv.Itoa(d.counter)
	name1 := HOST_VETH_PREFIX + id
	name2 := "veth-temp-" + id
	err := createVeth(name1, name2)
	if err != nil {
		return nil, err
	}

	program, ok := d.programMap[request.NetworkID]
	if !ok {
		return nil, errors.New("couldn't find network")
	}
	err = program.Attach(name1, name1+id)
	if err != nil {
		return nil, err
	}

	cid := connectionID{NetworkID: request.NetworkID, EndpointID: request.EndpointID}
	info := d.connectionInfoMap[cid]
	info.hostVethName = name1
	d.connectionInfoMap[cid] = info
	d.counter++
	return &network.JoinResponse{InterfaceName: network.InterfaceName{
		SrcName:   name2,
		DstPrefix: "eth",
	}}, nil
}

func (d *Driver) Leave(request *network.LeaveRequest) error {
	log.Printf("Leave(): NetworkID: %s, EndpointID: %s\n", request.NetworkID[:6], request.EndpointID[:6])
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

// input: 192.168.34.12/16
// output: 192.168.34.12
func removeSubnet(ipAddr string) string {
	return strings.Split(ipAddr, "/")[0]
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
