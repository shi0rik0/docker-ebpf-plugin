//go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#define likely(x)		__builtin_expect(!!(x), 1)
#define unlikely(x)		__builtin_expect(!!(x), 0)

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
	if (unlikely((void *)(eth_header + 1) > data_end)) {
		return TC_ACT_OK;
	}
	__u32 ip_addr;
	switch (eth_header->h_proto) {
	case bpf_htons(ETH_P_IP): {
		struct iphdr *ip_header = (struct iphdr *)(eth_header + 1);
		if (unlikely((void *)(ip_header + 1) > data_end)) {
			return TC_ACT_OK;
		}
		ip_addr = ip_header->daddr;
		break;
	}
	case bpf_htons(ETH_P_ARP): {
		void *arp_payload = eth_header + 1;
		if (unlikely(arp_payload + 28 > data_end)) {
			return TC_ACT_OK;
		}
		ip_addr = *(__u32 *)(arp_payload + 24);
		break;
	}
	default:
		return TC_ACT_OK;
		break;
	}
	__u32 *ifindex = bpf_map_lookup_elem(&ip_ifindex_map, &ip_addr);
	if (likely(ifindex)) {
		return bpf_redirect(*ifindex, 0);
	}
	return TC_ACT_OK;
}

char __license[] SEC("license") = "GPL";
