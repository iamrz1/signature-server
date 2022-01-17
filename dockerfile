# builder image
FROM golang:latest AS builder

# copy source code
WORKDIR /source
COPY . .

# fetch dependencies
RUN go mod download

# build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARC=amd64 go build -a -installsuffix cgo -o app .


# base image
FROM alpine:latest

# Security related package
RUN apk --no-cache add ca-certificates

# copy the binary
COPY --from=builder /source/app /usr/local/bin/app

# run the binary
ENTRYPOINT [ "app" ]