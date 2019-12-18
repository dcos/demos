build:
	docker run --rm -v $(PWD):/usr/src/actor -w /usr/src/actor \
		golang:1.12.5 \
		go get -t -v github.com/Shopify/sarama && \
		GOOS=linux GOARCH=386 CGO_ENABLED=0 \
		go build -ldflags="-s -w" -o fraudDisplay-linux

