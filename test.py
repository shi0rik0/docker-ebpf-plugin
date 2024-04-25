#!/usr/bin/env python3

import subprocess
import argparse
import json
import re

# 创建解析器
parser = argparse.ArgumentParser()

# 添加参数
parser.add_argument('-d', choices=['bridge', 'overlay', 'ipvlan', 'macvlan', 'ebpf'], required=True, help='')
parser.add_argument('-t', choices=['iperf3', 'netperf'], required=True, help='')
parser.add_argument('-f', type=str, help='')

A = {
    'bridge': 'bridge',
    'overlay': 'overlay-net',
    'ipvlan': 'ipvlan-net',
    'macvlan': 'macvlan-net',
    'ebpf': 'ebpf-net',
}

# 解析命令行参数
args = parser.parse_args()

def run(cmd):
    return subprocess.run(cmd.split(' '), stdout=subprocess.PIPE).stdout.decode('utf-8')

def parse_iperf3_result(s):
    s = [i for i in s.split('\n') if 'receiver' in i][0]
    match = re.search(r'(\d+)\s+Mbits/sec', s)
    if match:
        return match.group(1)
    else:
        return None

def parse_netperf_result(s):
    return s.split('\n')[-3].split()[-1]
    
def get_ip(driver):
    r = [None] * 2
    d = json.loads(run(f'docker network inspect {A[driver]}'))
    for i in d[0]['Containers'].values():
        if i['Name'] in ['container1', 'container2']:
            ip = i['IPv4Address'].split('/')[0]
            if i['Name'] == 'container1':
                r[0] = ip
            else:
                r[1] = ip
    return r


ip = get_ip(args.d)[0]
if args.t == 'iperf3':
    r = run(f'docker exec container2 iperf3 -c {ip} -f m')
    print(parse_iperf3_result(r))
elif args.t == 'netperf':
    r = run(f'docker exec container2 netperf -t TCP_RR -H {ip}')
    print(parse_netperf_result(r))