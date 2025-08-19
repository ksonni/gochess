build:
	cd app/cmd && go build -o gochess

run:
	make build && ./app/cmd/gochess

test:
	cd app && go test ./... -v
