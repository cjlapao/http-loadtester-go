############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /go/src/cjlapao/http-loadtester-go

COPY . .

# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/http_loadtester

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable.
COPY --from=builder /go/bin/http_loadtester /go/bin/http_loadtester
RUN chmod +x /go/bin/http_loadtester

EXPOSE 10000
ENTRYPOINT ["/go/bin/http_loadtester", "api"]
