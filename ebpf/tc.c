// go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
// #include <bpf/bpf_tracing.h>

#define TC_ACT_OK 0
#define ETH_P_IP bpf_htons(0x0800) /* Internet Protocol packet	*/
#define ETH_P_ARP bpf_htons(0x0806)

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__type(key, __u32);
	__type(value, __u32);
	__uint(max_entries, 4096);
} ip_ifindex_map SEC(".maps");

static __always_inline __u32 get_ip_addr(struct ethhdr *eth_header)
{
	// void *data = (void *)(uintptr_t)ctx->data;
	// void *data_end = (void *)(uintptr_t)ctx->data_end;
}

SEC("tc")
int tc_ingress(struct __sk_buff *ctx)
{
	if (ctx->data_end - ctx->data < 14) {
		return TC_ACT_OK;
	}
	struct ethhdr *eth_header = (struct ethhdr *)(uintptr_t)ctx->data;
	switch (eth_header->h_proto) {
	case bpf_htons(ETH_P_IP):
		break;
	case bpf_htons(ETH_P_ARP):
		break;
	default:
		break;
	}
	// bpf_printk("%u", ctx->data_end - ctx->data);
	// bpf_printk("%u", ctx->protocol);
	return TC_ACT_OK;
}

char __license[] SEC("license") = "GPL";
