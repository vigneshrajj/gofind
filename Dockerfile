FROM golang:1.23.2 AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/gofind/main.go

FROM alpine:latest
COPY --from=builder /main /main
COPY --from=builder /app/static /static

RUN ls -l /
EXPOSE 3005
ENTRYPOINT ["/main"]
