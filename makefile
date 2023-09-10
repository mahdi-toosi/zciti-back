ms:
	make migrate
	make seed

migrate:
	go run cmd/example/main.go --migrate

seed:
	go run cmd/example/main.go --seed

swag:
	swag init -g ./cmd/example/main.go