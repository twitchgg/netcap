#!/bin/sh
sudo ../build/net-cap-v1.0.0-linux-amd64 \
        --ngrep-path="/usr/bin/ngrep" \
        --bind-addr="127.0.0.1:8081" \
        --logger-level="debug"


