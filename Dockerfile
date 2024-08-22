FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server cmd/books/main.go

FROM golang:alpine
COPY --from=builder /app/server .
CMD ["./server"]