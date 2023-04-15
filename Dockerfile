# build stage
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# production stage
FROM alpine

COPY --from=builder /app/main /main

ENTRYPOINT ["/main", "app:serve"]
