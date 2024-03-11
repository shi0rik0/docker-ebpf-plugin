#!/bin/bash

for i in 'bridge' 'overlay' 'ipvlan' 'macvlan' 'ebpf'; do
    echo "$i"
    ./test.py -d "$i" -t iperf3
done