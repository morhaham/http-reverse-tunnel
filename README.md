# Instructions

## Creating TLS  self signed certificate
`mkdir ./cmd/tunnel-server/tls/ && cd ./cmd/tunnel-server/tls/`

`openssl genrsa -out server.key 2048`

`openssl genrsa -out server.key 2048`

`openssl ecparam -genkey -name secp384r1 -out server.key`

`openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`
 
