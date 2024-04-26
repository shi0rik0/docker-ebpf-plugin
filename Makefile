bin/denp: ebpf/vmlinux.h ebpf/tc_bpfeb.o XXX
	go build -o bin/denp .

ebpf/tc_bpfeb.o: ebpf/tc.c
	go generate ./ebpf/

ebpf/vmlinux.h:
	bpftool btf dump file /sys/kernel/btf/vmlinux format c > ebpf/vmlinux.h

run: bin/denp
	sudo bin/denp

clean:
	rm -rf bin ebpf/tc_bpfeb.go ebpf/tc_bpfeb.o ebpf/tc_bpfel.go ebpf/tc_bpfel.o ebpf/vmlinux.h

XXX: ;
