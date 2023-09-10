build:
	go build -o gochess

run:
	make build && ./gochess

test:
	go test ./...
