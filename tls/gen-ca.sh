#!/bin/bash
set -ex
pass=$1

echo $1

# generate CA's  key
openssl genrsa -aes256 -passout pass:"$pass" -out key 4096
openssl rsa -passin pass:"$pass" -in key -out ca.key

openssl req -config ca.cnf -key ca.key -new -x509 -days 7300 -sha256 -extensions v3_ca -out ca.crt