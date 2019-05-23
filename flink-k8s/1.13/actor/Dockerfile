FROM golang:1.12.5
ADD . /src
RUN set -x && \
    cd /src && \
    go get -t -v github.com/Shopify/sarama && \
    CGO_ENABLED=0 GOOS=linux go build -a -o fraudDisplay-linux

FROM scratch
COPY --from=0 /src/fraudDisplay-linux /
CMD ["/fraudDisplay-linux"]
