FROM golang:1.25 AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p bin && \
    go build -ldflags "-X main.Version=0.1.0" -o bin/app ./cmd/my_knowledges

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates netbase \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000

VOLUME /data/conf
VOLUME /app/uploads

CMD ["./app", "-conf", "/data/conf"]
