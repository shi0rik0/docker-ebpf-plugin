#!/usr/bin/env bash

docker rm -f container1
docker rm -f container2
docker network rm overlay-net
docker network rm ipvlan-net
docker network rm macvlan-net
docker network rm ebpf-net
docker swarm leave --force
sudo pkill denp