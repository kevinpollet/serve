FROM golang:1.17-alpine3.14 AS builder
WORKDIR /go/src/serve
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /go/bin/serve .

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/serve /usr/local/bin
USER 65534
EXPOSE 8080
ENTRYPOINT [ "serve", "--addr", ":8080" ]
CMD [ "--dir", "/var/www" ]
