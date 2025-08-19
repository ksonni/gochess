build:
	cd app && go build -o app ./cmd/

run:
	make build && ./app/app

test:
	cd app && go test ./... -v
