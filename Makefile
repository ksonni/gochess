build:
	cd app && go build -o app ./cmd/

up:
	docker-compose build && docker-compose up

test:
	cd app && go test ./... -v
