#!/bin/bash

ip netns add vm1
ip netns add vm2
ip link add veth1 type veth peer name veth2
ip link add veth3 type veth peer name veth4
ip link set veth1 up
ip link set veth2 up
ip link set veth3 up
ip link set veth4 up
ip addr add 10.0.0.2/24 dev veth1
ip addr add 10.0.0.3/24 dev veth3
ip link set veth1 netns vm1
ip link set veth3 netns vm2