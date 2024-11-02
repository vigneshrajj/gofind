FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /main ./cmd/gofind/main.go

EXPOSE 3005
CMD ["/main"]
