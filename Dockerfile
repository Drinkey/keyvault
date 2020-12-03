FROM golang:1.15.5

WORKDIR /go/src/keyvault
COPY ./keyvault .

ENV GOPROXY="https://mirrors.aliyun.com/goproxy/"

RUN go get -d -v ./...
RUN go install -v ./...

ENV CERT_CONF_DIR=/go/src/keyvault/etc
ENV CERT_DIR=/go/src/keyvault/etc
ENV DB_PATH=/go/src/keyvault/etc/vault.db

CMD ["keyvault"]
