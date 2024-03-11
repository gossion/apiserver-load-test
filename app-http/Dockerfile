# Start from a small Golang base image  
FROM golang:1.21-alpine as builder  
  
# Set the Current Working Directory inside the container  
WORKDIR /app  
  
# Copy go mod and sum files  
COPY go.mod go.sum ./
  
# Download all dependencies  
RUN go mod download  
  
# Copy the source from the current directory to the Working Directory inside the container  
COPY . .  
  
# Build the Go app  
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o apiserver-load-tester .  
  
# Start a new stage from scratch  
FROM alpine:latest    
  
WORKDIR /root/  
  
# Copy the Pre-built binary file from the previous stageCompletions call failed. Please try again.
# Copy the Pre-built binary file from the previous stage  
COPY --from=builder /app/apiserver-load-tester .  
  
# Command to run the executable  
CMD ["./apiserver-load-tester"]  
