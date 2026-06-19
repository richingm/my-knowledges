FROM golang:1.25 AS go-builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p bin && \
    go build -ldflags "-X main.Version=0.1.0" -o bin/app ./cmd/my_knowledges

FROM node:22 AS web-builder

WORKDIR /web

COPY web/package.json web/package-lock.json ./
RUN npm install

COPY web/ .

RUN npm run build

FROM nginx:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=go-builder /app/bin/app /usr/local/bin/app

COPY --from=web-builder /web/dist /usr/share/nginx/html

COPY web/nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /data/conf /app/uploads

EXPOSE 80 8000 9000

VOLUME /data/conf
VOLUME /app/uploads

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]