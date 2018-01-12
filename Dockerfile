FROM golang:alpine as builder

WORKDIR /go/src/github.com/mccraymt/ms-black-history
ADD . /go/src/github.com/mccraymt/ms-black-history/

RUN go build -o ms-geo-data .

FROM alpine:latest
RUN apk add --no-cache bash

WORKDIR /root/
COPY --from=builder /go/src/github.com/mccraymt/ms-black-history/ms-geo-data .
COPY --from=builder /go/src/github.com/mccraymt/ms-black-history/config.json .

CMD /root/ms-geo-data

EXPOSE 4000
