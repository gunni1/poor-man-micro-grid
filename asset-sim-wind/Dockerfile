# Build stage
FROM golang:1.25 AS builder

ENV GOOS=linux
ENV GOARCH=arm64
ENV CGO_ENABLED=0

WORKDIR /app
COPY shared/ ./shared/
COPY asset-sim-wind/ ./asset-sim-wind/
WORKDIR /app/asset-sim-wind
RUN go mod download
RUN go build -o /app/asset-sim-wind .

# Runtime
FROM gcr.io/distroless/static-debian12
WORKDIR /app/
COPY --from=builder /app/asset-sim-wind .
USER nonroot:nonroot
ENTRYPOINT ["/app/asset-sim-wind"]
