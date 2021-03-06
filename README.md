# keyvault

![Go](https://github.com/Drinkey/keyvault/workflows/Go/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/Drinkey/keyvault/branch/main/graph/badge.svg?token=DKOXF0NYVN)](https://codecov.io/gh/Drinkey/keyvault)


Practice Project: key vault to store your secrets, provide RESTFul API, build with Gin

> Just a practice project, it works in internal network. It's not a production-ready application.

Structure
```
+------------------+
|        API       |
+------------------+
|      Services    |
+------------------+
|  CertIO |   pkg  |
+---------+--------+
|       Models     |
+------------------+

```
# Features

- client and server use certificate to validate each other
- server perform as a CA, and able to issue certificate to client
- only authorized clients can communicate with keyvault
- the certificate must have OU property, server will use it as the namespace of secret to query database
- server has a master key to encrypt/decrypt the secrets using AES
- the master key is randomly generated and stored in a file that only visible to the user who start the service
- the master key is also encrypted in the database and only visible during runtime
- server has RESTFul API to store the secrets, get secrets, issue certificates
- communication is protected by TLS

The secret still visible as a plain text to authorized client. If the client choose to print it out or log it, there is nothing we can do. The sensitive data should never been seen anywhere.

# Usage

API doc can be found in /docs
# Test

## Start the server

```
$ cd <the_path_of_project>
$ export DB_PATH=/tmp/vault.db
$ cd keyvault
## Start keyvault server
$ go run keyvault.go
```
## Prepare certificates

**Server side**

CA certificate and privatekey, keyvault service certificate and private key are automatically generated when starting the server.

If you want to re-generate all certificates, just remove `keyvault/etc/ca.crt` and start the server again. You can use `keyvault/etc/ca.key` to sign client certificate request. The keyvault service will provide API to sign the request later.

**Client side**

Client need to generate private key manually. The following command will generate private key with RSA algorithm and key length is 4096.
```
$ openssl genrsa -out client.key 4096
```

Then client should create a CSR(Certificate Signing Request) file with OU specified, keyvault service will use OU to determine if client is authorized to access specific secret. For example, we set OU to `KUBERNETES`, and set other field to fit your need.
```
$ openssl req -new -nodes -key client.key -out client.csr -subj /C=CN/ST=SC/L=CD/O="KeyVault Client"/OU=KUBERNETES/CN=K8S.keyvault.org
```

Then use the CA key pairs to sign this request. (TODO: keyvault will provide API to do this, and this section should be updated)
```
$ openssl x509 -req -in client.csr -CA keyvault/etc/ca.crt  -CAkey keyvault/etc/ca_priv.key  -CAcreateserial -out client.crt
```

## Access API

Create a new namespace. The namespace value must be exactly the same with OU in `client.crt`
```
curl -L --key certs/client.key --cert certs/client.crt --cacert certs/ca.crt https://keyvault.org/v1/vault/ -X POST -d '{"name": "KUBERNETES"}'
```

Create a new secret under the namespace
```
curl --key certs/client.key --cert certs/client.crt --cacert certs/ca.crt https://keyvault.org/v1/vault/KUBERNETES -X POST -d '{"key": "admin_user", "value": "some-password"}'
```

Get the secret
```
curl --key certs/client.key --cert certs/client.crt --cacert certs/ca.crt https://keyvault.org/v1/vault/KUBERNETES\?q\=admin_user
```

# Development

go 1.15.5

## Build the Docker Image

```
$ git clone 
$ cd keyvault
$ docker build -t keyvault:latest .
# or
$ docker-compose build
```
## Start the service

### First look

Start the service and listen on local port 443 in foreground. Ctl+c will terminate the service, and remove the container.
```
docker run --rm -p 443:443 --name keyvault keyvault:latest
```

### Data persistance

Persist certificate and database file.

```
$ export HOST_DATA_DIR=/some/where/on/host
$ docker-compose up -d
```

## Build binary with Docker

```
$ docker run --rm \
    -v "$PWD/keyvault":/usr/src/keyvault \
    -w /usr/src/keyvault \
    -e GOOS=linux \
    -e GOARCH=amd64 \
    -e CGO_ENABLED=1 \
    keyvault:debug \
    go build -v -o keyvault-linux-amd64
```

Note, build with `GOOS=darwin` and `CGO_ENABLED=1` will fail on my macOS.

# Other usage

This tool can be also used for generate CA and web certificate files.

Write your customized configuration file, update the `certificates` part with certificates details.
```sh
$ mkdir psonocdc
# specify the configuration file path
$ export KV_CONFIG_FILE=$PWD/psono.json
# specify the dir path to store generated certificates
$ export KV_CERT_DIR=$PWD/psonocdc
# start the service then certificates will be generated automatically
$ go run keyvaultd.go
...
# checkout the generated certificates
$ ls psonocdc
ca.crt        ca_priv.key   cert.pem      cert_priv.key
```
