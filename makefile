dms:
	make drop
	make migrate
	make seed

drop:
	go run cmd/example/main.go --drop-all-tables

migrate:
	go run cmd/example/main.go --migrate

seed:
	go run cmd/example/main.go --seed

swag:
	swag init -g ./cmd/example/main.go --outputTypes "json"