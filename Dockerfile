FROM golang:1.18.0-alpine3.15 as builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod .
COPY go.sum .

COPY . .

# Build the Go app
RUN go build  -o ./hdrcrypto ./cmd/hdrcrypto

# generate clean, final image for end users
FROM alpine:3.15
WORKDIR /app

COPY --from=builder /build/hdrcrypto .

# This container exposes port 8080 to the outside world
EXPOSE 3000

ENTRYPOINT ["/app/hdrcrypto"]
CMD ["serve"]