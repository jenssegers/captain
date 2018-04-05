PWD=$(shell pwd)

build:
	docker run -it --init --rm -v $(PWD):/code -w /code instrumentisto/glide install
	docker run -it --init --rm -v $(PWD):/go/src/captain -e GOPATH=/go -w /go/src/captain -e GOOS=darwin golang go build -o bin/osx/captain
	docker run -it --init --rm -v $(PWD):/go/src/captain -e GOPATH=/go -w /go/src/captain -e GOOS=linux golang go build -o bin/linux/captain
