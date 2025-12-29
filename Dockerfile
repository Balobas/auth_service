FROM --platform=$BUILDPLATFORM golang:1.24.11-alpine AS builder

COPY . /github.com/balobasta/auth_service/src/
WORKDIR /github.com/balobasta/auth_service/src/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/auth_service cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/balobasta/auth_service/src/bin/auth_service .
COPY --from=builder /github.com/balobasta/auth_service/src/local.env .

EXPOSE 50051

ENTRYPOINT ["./auth_service", "-config-path=local.env"]