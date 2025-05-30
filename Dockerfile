FROM golang:1.22

WORKDIR /app

# Copy the application code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

EXPOSE 8080

# Entry point: wait for the database and then start the Go app
CMD ["./main"]
