rm *.pem

# Generate CA's private key & self-signed cert
openssl req -x509 -newkey rsa:4096 -nodes -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=HU/ST=Budapest/L=Budapest/O=Petifies/OU=Petifies/CN=petifies.com/emailAddress=vantoan.hk14@gmail.com"

openssl x509 -in ca-cert.pem -noout -text

# Generate server's private key and certificate signing request
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=HU/ST=Budapest/L=Budapest/O=Petifies/OU=Petifies/CN=petifies.com/emailAddress=vantoan.hk14@gmail.com"

# Generate server
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

openssl x509 -in server-cert.pem -noout -text