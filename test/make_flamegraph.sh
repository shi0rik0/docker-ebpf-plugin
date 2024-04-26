#!/usr/bin/env bash

sudo perf record -F 99 -e cpu-clock -ag -- sleep "$1"
sudo perf script > out.perf
FlameGraph/stackcollapse-perf.pl out.perf > out.folded
FlameGraph/flamegraph.pl out.folded > "$2"

rm -f perf.data perf.data.old out.perf out.folded