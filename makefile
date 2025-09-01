rebuild:
	make drop
	make deleteFilesInStorage
	make migrate
	make generateNecessaryData
	make seed

drop:
	go run cmd/seeder/seeder.go --drop-all-tables

deleteFilesInStorage:
	go run cmd/seeder/seeder.go --delete-files-in-storage

migrate:
	go run cmd/seeder/seeder.go --migrate

generateNecessaryData:
	go run cmd/seeder/seeder.go --generate-necessary-data

seed:
	go run cmd/seeder/seeder.go --seed

swag:
	swag init -g ./cmd/main/main.go --outputTypes "json"

release:
	GOPROXY="https://goproxy.io,direct" go build -mod=vendor -o ./tmp/main ./cmd/main/main.go