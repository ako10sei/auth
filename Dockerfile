FROM golang:1.22.7-alpine AS builder

COPY . /github.com/ako10sei/auth/source
WORKDIR /github.com/ako10sei/auth/source

RUN go mod download
RUN go build -o ./bin/crud cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/ako10sei/auth/source/bin/crud .

CMD ["./crud"]