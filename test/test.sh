#!/usr/bin/env bash

./setup_test.sh
mkdir -p graphs

for i in 'bridge' 'overlay' 'ebpf'; do
    echo "Testing the $i network by $1."
    ./_test.py -d "$i" -t "$1" -f "graphs/${i}_$1.svg"
done

./clean_test.sh
