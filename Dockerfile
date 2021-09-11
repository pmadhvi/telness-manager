# Using golang alpine image
FROM golang:1.15-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum .env ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the go app
RUN go build -o telness-manager cmd/main.go

# expose port 8080 from container
EXPOSE 8080

CMD ["./telness-manager"]