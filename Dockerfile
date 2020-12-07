FROM golang:1.15.5

WORKDIR /go/src/keyvault
COPY ./keyvault .

ENV GOPROXY="https://mirrors.aliyun.com/goproxy/"

RUN go get -d -v ./...
RUN go install -v ./...

ENV CERT_CONF_FILE=/go/src/keyvault/cert.json
ENV CERT_DIR=/go/src/keyvault/certs
ENV DB_PATH=/go/src/keyvault/vault.db

CMD ["keyvault"]
