#!/bin/bash

openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=kubexdp-webhook-ca" -days 365 -out ca.crt

# Generate webhook server certificate
openssl genrsa -out webhook.key 2048
openssl req -new -key webhook.key -subj "/CN=kubexdp-webhook-service.kubexdp-system.svc" -out webhook.csr
openssl x509 -req -in webhook.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out webhook.crt -days 365 -extensions v3_req

# Create secret manually
kubectl create secret tls webhook-server-cert \
    --cert=webhook.crt \
    --key=webhook.key \
    -n kubexdp-system

