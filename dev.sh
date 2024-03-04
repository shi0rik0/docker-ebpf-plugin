#!/bin/bash

go generate ./ebpf/
sudo go run ./ebpf/