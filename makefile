rebuild:
	make drop
	make deleteFilesInStorage
	make migrate
	make generateNecessaryData
	make seed

drop:
	go run cmd/example/seeder.go --drop-all-tables

deleteFilesInStorage:
	go run cmd/example/seeder.go --delete-files-in-storage

migrate:
	go run cmd/example/seeder.go --migrate

generateNecessaryData:
	go run cmd/example/seeder.go --generate-necessary-data

seed:
	go run cmd/example/seeder.go --seed

swag:
	swag init -g ./cmd/example/main.go --outputTypes "json"

release:
	go build -mod=vendor -o ./tmp/main ./cmd/example/main.go