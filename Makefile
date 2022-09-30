PWD=$(shell pwd)

build:
	docker run -it --init --rm -v $(PWD):/go/src/captain -w /go/src/captain golang go mod tidy
	docker run -it --init --rm -v $(PWD):/go/src/captain -e GOPATH=/go -w /go/src/captain -e GOOS=darwin golang go build -o dist/captain-osx
	docker run -it --init --rm -v $(PWD):/go/src/captain -e GOPATH=/go -w /go/src/captain -e GOOS=linux golang go build -o dist/captain-linux
	docker run -it --init --rm -v $(PWD):/go/src/captain -e GOPATH=/go -w /go/src/captain -e GOOS=windows golang go build -o dist/captain.exe
