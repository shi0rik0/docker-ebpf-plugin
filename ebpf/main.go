package main

import (
	"log"
	"os"
	"os/signal"

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
		log.Fatal("Loading eBPF objects:", err)
	}
	defer objs.Close()

	iface, err := netlink.LinkByIndex(225)
	if err != nil {
		log.Fatal(err)
	}

	// Attach count_packets to the network interface.
	err = tc.AttachTC(objs.TcIngress, iface, tc.INGRESS, "tc_ingress")
	if err != nil {
		log.Fatal("Attaching TC:", err)
	}

	log.Printf("AttachTC成功！\n")

	// Periodically fetch the packet counter from PktCount,
	// exit the program when interrupted.
	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)
	for {
		select {
		case <-stop:
			log.Print("Received signal, exiting..")
			err := tc.DetachTC(objs.TcIngress, iface, tc.INGRESS)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}
