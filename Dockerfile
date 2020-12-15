FROM golang:alpine3.12 as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
# at first download all dependencies, so this step can be omitted on next run
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /postman-mockserver
FROM alpine:3.12
COPY --from=builder /postman-mockserver /app/postman-mockserver
COPY /docker-entrypoint.sh /docker-entrypoint.sh
RUN ["chmod", "+x", "/docker-entrypoint.sh"]
COPY ./config.yaml /app/config/config.yaml
ENTRYPOINT ["./docker-entrypoint.sh"]