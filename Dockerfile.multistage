FROM golang:1.12 as builder

LABEL maintainer="pashaborisyk <pashaborisyk@gmail.com>"

WORKDIR /go/src/look-and-like-search-converter

COPY . .

RUN go get -d -v ./...

WORKDIR /go/src/look-and-like-search-converter/main

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/look-and-like-search-converter .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/bin/look-and-like-search-converter .

CMD ["./look-and-like-search-converter"]