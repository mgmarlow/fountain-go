BINARY_NAME=ftn

build:
	go build -o ${BINARY_NAME} main.go lexer.go parser.go

run:
	go run main.go

test:
	go test

benchmark:
	go test -bench=.

clean:
	go clean
	rm ${BINARY_NAME}
