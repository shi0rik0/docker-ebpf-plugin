docker network disconnect ebpfnet1 test1
docker network disconnect ebpfnet1 test2
docker rm -f test1
docker rm -f test2
docker network rm ebpfnet1