FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/cespare/reflex@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -gcflags="all=-N -l" -o main ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

COPY --from=builder /app/main .

EXPOSE 8080 2345

CMD ["dlv", "exec", "./main", "--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient", "--log"]
CMD ["reflex", "-r", "\\.go$", "-s", "--", "go", "run", "./cmd/api"]