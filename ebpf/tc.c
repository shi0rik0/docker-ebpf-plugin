//go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#define TC_ACT_OK 0
#define ETH_P_IP 0x0800
#define ETH_P_ARP 0x0806

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__type(key, __be32);
	__type(value, __u32);
	__uint(max_entries, 4096);
} ip_ifindex_map SEC(".maps");

SEC("tc")
int tc_ingress(struct __sk_buff *ctx)
{
	if (ctx->data_end - ctx->data < 14) {
		return TC_ACT_OK;
	}
	struct ethhdr *eth_header = (struct ethhdr *)(uintptr_t)ctx->data;

	__be32 ip_addr;
	switch (eth_header->h_proto) {
	case bpf_htons(ETH_P_IP): {
		struct iphdr *ip_header = (struct iphdr *)(eth_header + 1);
		ip_addr = ip_header->daddr;
		break;
	}
	case bpf_htons(ETH_P_ARP): {
		struct arphdr *arp_header = (struct arphdr *)(eth_header + 1);
		void *arp_payload = eth_header + 1;
		ip_addr = *(__be32 *)(arp_payload + 16);
		break;
	}
	default:
		return TC_ACT_OK;
		break;
	}
	bpf_printk("%u", bpf_ntohl(ip_addr));
	return TC_ACT_OK;
}

char __license[] SEC("license") = "GPL";
