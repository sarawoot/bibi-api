# Use the official Golang image to create a build artifact.
# https://hub.docker.com/_/golang
FROM golang:1.18.5-alpine3.15 AS builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
# -- This shouldn't run because we already have vendor dir
# COPY go.* ./
# RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
# -mod=vendor to not re-download vendor packages
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -v -o server ./cmd/server/main.go

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3.15.0
RUN apk add --no-cache ca-certificates=20220614-r0
RUN apk add --no-cache tzdata=2022c-r0

WORKDIR /app
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server server

# Go1.15 allow to put timezone data by import in main program
# ENV TZ=Asia/Bangkok
ENV SERVER_ADDR=0.0.0.0:80

# Run the web service on container startup.
CMD ["./server"]


