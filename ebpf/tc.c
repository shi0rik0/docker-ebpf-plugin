//go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#define TC_ACT_OK 0
#define ETH_P_IP 0x0800
#define ETH_P_ARP 0x0806

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__type(key, __u32);
	__type(value, __u32);
	__uint(max_entries, 4096);
} ip_ifindex_map SEC(".maps");

SEC("tc")
int tc_ingress(struct __sk_buff *ctx)
{
	void *data_end = (void *)(uintptr_t)ctx->data_end;
	struct ethhdr *eth_header = (struct ethhdr *)(uintptr_t)ctx->data;
	if ((void *)(eth_header + 1) > data_end) {
		return TC_ACT_OK;
	}
	__u32 ip_addr;
	switch (eth_header->h_proto) {
	case bpf_htons(ETH_P_IP): {
		bpf_printk("ip packet");
		struct iphdr *ip_header = (struct iphdr *)(eth_header + 1);
		if ((void *)(ip_header + 1) > data_end) {
			return TC_ACT_OK;
		}
		ip_addr = ip_header->daddr;
		break;
	}
	case bpf_htons(ETH_P_ARP): {
		bpf_printk("arp packet");
		void *arp_payload = eth_header + 1;
		if (arp_payload + 28 > data_end) {
			bpf_printk("%d", data_end - arp_payload);
			return TC_ACT_OK;
		}
		ip_addr = *(__u32 *)(arp_payload + 24);
		break;
	}
	default:
		bpf_printk("unknown packet");
		return TC_ACT_OK;
		break;
	}
	bpf_printk("lookup for %u", ip_addr);
	__u32 *ifindex = bpf_map_lookup_elem(&ip_ifindex_map, &ip_addr);
	if (ifindex) {
		bpf_printk("redirect to %u", *ifindex);
		return bpf_redirect(*ifindex, 0);
	}
	bpf_printk("lookup failed", *ifindex);
	return TC_ACT_OK;
}

char __license[] SEC("license") = "GPL";
