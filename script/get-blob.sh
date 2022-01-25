#!/bin/sh

ref="${1:-library/busybox:latest}"
repo="${ref%:*}"
tag="${ref##*:}"
digest="$2"
token=$(curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:${repo}:pull" \
        | jq -r '.token')
curl -H "Authorization: Bearer $token" \
     -k \
     --proxy https://localhost:9443 \
     --proxy-insecure \
     -s -L -o - "https://registry-1.docker.io/v2/${repo}/blobs/${digest}"
