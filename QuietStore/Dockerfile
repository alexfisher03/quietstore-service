# compile binary
FROM golang:1.25.0-alpine AS builder
WORKDIR /app
COPY QuietStore/go.mod QuietStore/go.sum ./
RUN go mod download

COPY QuietStore/. .
RUN go build -o server ./cmd/server

# run stage (linux)
FROM alpine:3.18
# https
RUN apk add --no-cache ca-certificates 
WORKDIR /app
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["./server"]
