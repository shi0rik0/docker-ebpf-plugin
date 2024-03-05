package ebpf

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/shi0rik0/docker-ebpf-plugin/tc"
	"github.com/vishvananda/netlink"
)

func main() {
	// Remove resource limits for kernels <5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Removing memlock:", err)
	}

	// Load the compiled eBPF ELF and load it into the kernel.
	var objs tcObjects
	if err := loadTcObjects(&objs, nil); err != nil {
		var ve *ebpf.VerifierError
		if errors.As(err, &ve) {
			for _, i := range ve.Log {
				fmt.Print(i + "\n")
			}
		}
		log.Fatal("Loading eBPF objects:", err)
	}
	defer objs.Close()

	iface1, err := netlink.LinkByIndex(11)
	if err != nil {
		log.Fatal(err)
	}
	iface2, err := netlink.LinkByIndex(13)
	if err != nil {
		log.Fatal(err)
	}

	// Attach count_packets to the network interface.
	err = tc.AttachTC(objs.TcIngress, iface1, tc.INGRESS, "tc_ingress1")
	if err != nil {
		log.Fatal("Attaching TC:", err)
	}
	err = tc.AttachTC(objs.TcIngress, iface2, tc.INGRESS, "tc_ingress2")
	if err != nil {
		log.Fatal("Attaching TC:", err)
	}

	log.Printf("AttachTC成功！\n")

	ip1, _ := ipAddrToUint32("10.0.0.2")
	log.Printf("ip1: %d", ip1)
	ip2, _ := ipAddrToUint32("10.0.0.3")
	log.Printf("ip2: %d", ip2)

	err = objs.IpIfindexMap.Put(ip1, uint32(11))
	if err != nil {
		log.Print(err)
	}
	err = objs.IpIfindexMap.Put(ip2, uint32(13))
	if err != nil {
		log.Print(err)
	}

	// Periodically fetch the packet counter from PktCount,
	// exit the program when interrupted.
	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)
	for {
		select {
		case <-stop:
			log.Print("Received signal, exiting..")
			err := tc.DetachTC(objs.TcIngress, iface1, tc.INGRESS)
			if err != nil {
				log.Fatal(err)
			}
			err = tc.DetachTC(objs.TcIngress, iface2, tc.INGRESS)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}
