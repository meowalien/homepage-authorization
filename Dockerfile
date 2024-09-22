FROM golang:1.22.4-alpine3.20 as builder

WORKDIR /app

COPY . .
#RUN go mod download

RUN go build -a -mod vendor -o /tmp/server ./main.go

FROM alpine:3.20 as runtime

WORKDIR /app

COPY --from=builder /tmp/server ./server

EXPOSE 8080

ENV GIN_MODE=release

CMD ["/app/server"]