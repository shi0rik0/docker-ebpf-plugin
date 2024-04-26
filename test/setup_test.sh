#!/usr/bin/env bash

sudo ../bin/denp &
docker swarm init
docker network create --driver overlay --attachable overlay-net
docker network create --driver ipvlan ipvlan-net
docker network create --driver macvlan macvlan-net
docker network create --driver ebpf ebpf-net
docker run -d --name container1 cilium/netperf
docker run -d --name container2 cilium/netperf
docker exec -d container1 iperf3 -s
for i in 'overlay-net' 'ipvlan-net' 'macvlan-net' 'ebpf-net'; do
    for j in 'container1' 'container2'; do
        docker network connect "$i" "$j"
    done
done
