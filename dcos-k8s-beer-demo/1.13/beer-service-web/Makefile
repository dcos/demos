build:
	docker run --rm -v $(PWD):/usr/src/beer -w /usr/src/beer \
		golang:1.12.5 \
		go build -ldflags="-s -w" -o beer-service-web-linux

