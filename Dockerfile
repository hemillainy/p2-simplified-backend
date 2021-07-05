FROM golang:1.16-alpine3.13 AS builder

WORKDIR /app
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ./


#FROM alpine:3.13
#WORKDIR /app
#COPY . .
#COPY --from=builder /app/migrate.linux-amd64 ./migrate

ENTRYPOINT ["sh", "scripts/init-db.sh"]