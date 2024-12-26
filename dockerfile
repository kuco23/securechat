FROM golang:1.23.2-bookworm

# Set the working directory for the project
WORKDIR /go/src/securechat

# Copy the go.mod and go.sum files first to enable dependency download
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application files
COPY ./securechat ./securechat
COPY ./home.html ./home.html

# Build the application binary and place it in /go/bin/
RUN go build -o /go/bin/securechat ./securechat

# Expose port 8080 for the service
EXPOSE 8080

# Define the command to run the application
CMD ["/go/bin/securechat"]