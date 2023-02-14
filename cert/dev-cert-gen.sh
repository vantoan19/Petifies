#!/bin/bash

rm ca.key ca.crt tls.key req.pem tls.crt ca.srl
# Generate CA's private key & self-signed cert
openssl req -x509 -newkey rsa:4096 -nodes -days 365 -keyout ca.key -out ca.crt -subj "/C=HU/ST=Budapest/L=Budapest/O=Petifies/OU=Petifies/CN=petifies.com/emailAddress=vantoan.hk14@gmail.com"

# Generate server's private key and certificate signing request
openssl req -newkey rsa:4096 -nodes -keyout tls.key -out req.pem -subj "/C=HU/ST=Budapest/L=Budapest/O=Petifies/OU=Petifies/CN=localhost/emailAddress=vantoan.hk14@gmail.com"

# Sign server's CSR with CA's private key and get back server cert
openssl x509 -req -in req.pem -days 60 -CA ca.crt -CAkey ca.key -CAcreateserial -out tls.crt -extfile ext.cnf