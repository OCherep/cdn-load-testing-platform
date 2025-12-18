FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o agent cmd/agent/main.go
ENTRYPOINT ["./agent"]
