FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

ENV CGO_ENABLED="0"
ENV GOOS="linux"

RUN go mod download
RUN go build -o goxy ./cmd/goxy/main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/goxy .

ENTRYPOINT [ "/app/goxy" ]
