FROM golang:latest AS builder
RUN apt-get update -y && apt-get upgrade -y
WORKDIR /go/src/
RUN git clone https://github.com/safronovD/spbstu-smart-library.git
WORKDIR /go/src/spbstu-smart-library/connector/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/connector ./src/*.go

FROM alpine:latest
RUN apk -U upgrade
WORKDIR /connector-app/
COPY --from=0 /go/bin/connector .
COPY --from=0 /go/src/spbstu-smart-library/connector/config.yaml .
