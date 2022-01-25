#!/bin/sh

ref="${1:-library/busybox:latest}"
repo="${ref%:*}"
tag="${ref##*:}"
api="application/vnd.docker.distribution.manifest.v2+json"
token=$(curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:${repo}:pull" \
        | jq -r '.token')
curl -H "Accept: ${api}" \
     -H "Authorization: Bearer $token" \
     -s "https://registry-1.docker.io/v2/${repo}/manifests/${tag}" | jq .
