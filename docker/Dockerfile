FROM golang:1.25-alpine AS builder

LABEL maintainer="netphantom"

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o homeSerBot .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /build/homeSerBot .
COPY --from=builder /build/ui ./ui
EXPOSE 4000

ENTRYPOINT ["./homeSerBot"]