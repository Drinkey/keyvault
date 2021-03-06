FROM golang:1.15.5

WORKDIR /go/src/keyvault
COPY . .

ENV GOPROXY="https://mirrors.aliyun.com/goproxy/"

RUN go get -d -v ./...
RUN go install -v ./...

ENV KV_CONFIG_FILE=/go/src/keyvault/keyvaultd-config.json
ENV KV_CERT_DIR=/go/src/keyvault/certs
ENV KV_DB_PATH=/go/src/keyvault/vault.db

CMD ["keyvaultd"]
