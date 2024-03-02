//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

#define TC_ACT_OK 0
#define ETH_P_IP  0x0800 /* Internet Protocol packet	*/

SEC("tc")
int tc_ingress(struct __sk_buff *ctx)
{
	return bpf_redirect(223, 0);
}

char __license[] SEC("license") = "GPL";
