#!/usr/bin/env bash

./setup_test.sh
mkdir -p graphs

for i in 'bridge' 'overlay' 'ebpf'; do
    echo "$i"
    ./test.py -d "$i" -t iperf3 -f "graphs/${i}_iperf3.svg"
done

./clean_test.sh
