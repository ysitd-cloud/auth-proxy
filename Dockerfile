FROM ysitd/dep AS builder

WORKDIR /go/src/app.ysitd/proxy

COPY . .

RUN dep ensure -v -vendor-only && \
    go install -v

FROM ysitd/binary

COPY --from=builder /go/bin/proxy /

CMD ["/proxy"]