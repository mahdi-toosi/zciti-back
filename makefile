rebuild:
	make drop
	make migrate
	make generateNecessaryData
	make seed

drop:
	go run cmd/example/main.go --drop-all-tables

migrate:
	go run cmd/example/main.go --migrate

generateNecessaryData:
	go run cmd/example/main.go --generate-necessary-data

seed:
	go run cmd/example/main.go --seed

swag:
	swag init -g ./cmd/example/main.go --outputTypes "json"