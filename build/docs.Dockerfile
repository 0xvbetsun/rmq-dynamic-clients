FROM golang:1.17-alpine AS builder

WORKDIR /build

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed 
# for our image and build the docs.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o doc-serv ./cmd/docs/main.go

FROM scratch

# Copy binary and config files from /build 
# to root folder of scratch container.
COPY --from=builder ["/build/doc-serv", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/doc-serv"]
