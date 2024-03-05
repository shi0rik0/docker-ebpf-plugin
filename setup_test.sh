docker network create --driver ebpf ebpfnet1
docker run -dit --name test1 cilium/netperf
docker run -dit --name test2 cilium/netperf
docker network connect ebpfnet1 test1
docker network connect ebpfnet1 test2
