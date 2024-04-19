rebuild:
	make drop
	make migrate
	make generateNecessaryData
	make seed

drop:
	go run cmd/example/seeder.go --drop-all-tables

migrate:
	go run cmd/example/seeder.go --migrate

generateNecessaryData:
	go run cmd/example/seeder.go --generate-necessary-data

seed:
	go run cmd/example/seeder.go --seed

swag:
	swag init -g ./cmd/example/main.go --outputTypes "json"