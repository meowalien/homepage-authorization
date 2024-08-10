FROM golang:1.22.4-alpine3.20 as builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o ./app

FROM alpine:3.20 as runtime

WORKDIR /app

COPY --from=builder /app/app ./app

EXPOSE 8080

ENV GIN_MODE=release

CMD ["/app/app"]