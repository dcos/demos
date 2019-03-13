build:
	docker run --rm -v $(PWD):/usr/src/generator -w /usr/src/generator \
		golang:1.12.0 \
		go get -t -v github.com/Shopify/sarama && \
		GOOS=linux GOARCH=386 CGO_ENABLED=0 \
		go build -ldflags="-s -w" -o generator-linux

