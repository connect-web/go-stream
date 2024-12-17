# Build stage
FROM golang:1.21 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o go-scraper-grpc ./cmd/server/main.go

# Run stage
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/go-scraper-grpc /go-scraper-grpc

EXPOSE 50051
ENTRYPOINT ["/go-scraper-grpc"]
