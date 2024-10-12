# Builder
FROM golang:1.23 AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build

# Final container
FROM scratch
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]
