FROM golang:1.17-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./

# Copy the code into the container.
COPY ./cmd/docs/main.go .

# Copy static assets
ADD ./web ./web

# Set necessary environment variables needed 
# for our image and build the docs.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o docs .

FROM scratch

# Copy binary and config files from /build 
# to root folder of scratch container.
COPY --from=builder ["/build/docs", "/"]
COPY --from=builder ["/build/web", "/web"]

# Command to run when starting the container.
ENTRYPOINT ["/docs"]
