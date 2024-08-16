FROM golang:1.21.4-alpine AS builder

COPY . /github.com/balobasta/auth_service_bln/src/
WORKDIR /github.com/balobasta/auth_service_bln/src/

RUN go build -o ./bin/auth_service cmd/server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/balobasta/auth_service_bln/src/bin/auth_service .
COPY --from=builder /github.com/balobasta/auth_service_bln/src/bin/.env .

ENTRYPOINT ["./auth_service", "-config-path=.env"]