FROM golang:1.15-alpine as builder

ENV GO111MODULE=on

ENV PORT=8000

WORKDIR /app
COPY . .

RUN go build -o keyval src/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/keyval .
CMD ["./keyval"]