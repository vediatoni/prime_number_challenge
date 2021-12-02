FROM golang:1.17.3-alpine as builder
WORKDIR /app

COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY ./cmd/input ./cmd/input
COPY ./config/dev.input.yaml ./config/dev.input.yaml
COPY ./internal/input ./internal/input
COPY ./pkg ./pkg

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
WORKDIR /app/cmd/input
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o /app/server

FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/config /config
COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]
