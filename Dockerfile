FROM golang:1.18.1 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM alpine:3.15
LABEL maintainer="Harvey"
ENV APP_ENV=prod
WORKDIR /app

EXPOSE 8080
EXPOSE 8081

COPY --from=builder /app/main /app/
COPY wsclient.html /app/

ENTRYPOINT ["/app/main"]