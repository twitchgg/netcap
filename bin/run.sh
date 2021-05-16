#!/bin/sh
sudo ../build/net-cap-v1.0.0-linux-amd64 \
        --dev="eth0" \
        --keyword="ali|baidu" \
        --ngrep-path="/usr/bin/ngrep" \
        --ports="52,53,54" \
        --host-ip="10.200.200.1" \
        --dump-path="./net_dump.pcap" \


