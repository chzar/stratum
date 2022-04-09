#!/bin/bash

digest="a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
token=$(curl \
    --silent \
    "https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/ubuntu:pull" \
    | jq -r '.token')
curl -H "Authorization: Bearer $token" \
     -k \
     --proxy https://0.0.0.0:9443 \
     --proxy-insecure \
     -s -L -o - "https://registry-1.docker.io/v2/library/ubuntu/blobs//sha256:${digest}" | sha256sum
