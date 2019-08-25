# STEP 1 - Build the binary

# use the golang image as build image
FROM golang

# activate go modules
ENV GO111MODULE on

# copy the local package files to the container's workspace.
COPY . /go/src/github.com/dj95/deception-proxy

# set the working directory to build the application
WORKDIR /go/src/github.com/dj95/deception-proxy

# compile the program
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -a \
    -installsuffix cgo \
    -o /go/bin/deception-proxy \
    /go/src/github.com/dj95/deception-proxy/cmd/deception-proxy/main.go


# STEP 2 - Build a minimal container

# start from scratch
FROM alpine

# install common ca-certificates
RUN apk update \
 && apk add --no-cache \
        ca-certificates \
 && rm -rf /var/cache/apk/*

# copy the static executable
COPY --from=0 /go/bin/deception-proxy /deception-proxy

# copy the config and web content into the container
COPY configs/config.yml /config.yml

# define the entrypoint
ENTRYPOINT ["/deception-proxy"]

# expose port 8080
EXPOSE 8080
