FROM golang:1.17-alpine3.14 AS builder
WORKDIR /go/src/srv
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /go/bin/srv cmd/srv/main.go

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/srv /usr/local/bin
USER 65534
EXPOSE 8080
ENTRYPOINT [ "srv", "--addr", ":8080" ]
CMD [ "--dir", "/var/www" ]
