FROM golang:alpine as builder

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build github.com/PaBah/GophKeeper/cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 3200

ENTRYPOINT ["/root/server", "-d", "host=postgres_db user=postgres password=postgres dbname=postgres", ""]