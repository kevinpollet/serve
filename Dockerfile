FROM golang:1.14 AS builder
WORKDIR /go/src/serge
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /go/bin/serge cmd/serge/main.go

FROM gcr.io/distroless/base
COPY --chown=65534:65534 --from=builder /go/bin/serge .
USER 65534
ENTRYPOINT [ "./serge" ]