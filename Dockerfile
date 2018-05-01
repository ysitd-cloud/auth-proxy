FROM ysitd/dep AS builder

WORKDIR /go/src/code.ysitd.cloud/proxy

COPY . .

RUN dep ensure -v -vendor-only && \
    go install -v

FROM ysitd/binary

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/proxy /

CMD ["/proxy"]