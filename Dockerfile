FROM golang:1.9.0 AS builder

ADD . /go/src/webcrawler/
WORKDIR /go/src/webcrawler/

RUN go get -v

RUN CGO_ENABLED=0 GOOS=linux go build -o server .


FROM alpine
WORKDIR /server

RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

RUN apk add -U tzdata
ADD config /server/config

COPY --from=builder /go/src/webcrawler/server .

CMD ["./server"]

EXPOSE 8080
