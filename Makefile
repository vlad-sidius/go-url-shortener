BINARY_NAME=shortener

all: build

clean:
	rm -rf dist/

build: clean lint
	go build -o dist/${BINARY_NAME} cmd/shortener/*.go

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...

report: coverage
	go tool cover -html=coverage.out

lint:
	go fmt ./...
	go vet ./...

autotest:
	./shortenertest -test.v -test.run=^TestIteration1$$ -binary-path=dist/${BINARY_NAME}
