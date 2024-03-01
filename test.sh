docker network create --driver ebpf ebpfnet1
docker run -dit --name test1 cilium/netperf
docker network connect ebpfnet1 test1
ip link show | grep "veth-"
docker exec test1 ip link show
docker exec test1 ip route show
docker network disconnect ebpfnet1 test1
docker exec test1 ip link show
docker exec test1 ip route show
docker rm -f test1
docker network rm ebpfnet1
