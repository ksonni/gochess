build:
	cd cmd && go build -o gochess

run:
	make build && ./cmd/gochess

test:
	go test ./... -v
