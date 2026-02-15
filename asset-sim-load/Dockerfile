# Build stage
FROM golang:1.25 AS builder

ENV GOOS=linux
ENV GOARCH=arm64
ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go mod tidy && go mod download
RUN go build -o asset-sim-load main.go

# Runtime
FROM gcr.io/distroless/static-debian12
WORKDIR /app/
COPY --from=builder /app/asset-sim-load .
USER nonroot:nonroot
ENTRYPOINT ["/app/asset-sim-load"]
