run:
	docker compose up --build

tests:
	go test -v ./...
