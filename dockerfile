############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /go/src/cjlapao/http-load-tester

COPY . .

WORKDIR /go/src/cjlapao/http-load-tester/src

# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/http-load-tester

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable.
COPY --from=builder /go/bin/http-load-tester /go/bin/http-load-tester

WORKDIR /go/bin

EXPOSE 80
ENTRYPOINT ["/go/bin/http-load-tester", "--api"]
