# Multi stage docker build. It will help in reducing the size of the final image. As in the finale image we only need to run the binary.

# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go 

# Run Stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]