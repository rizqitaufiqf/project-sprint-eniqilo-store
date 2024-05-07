build-dev:
	docker compose build

restart-dev:
	docker restart store-web

run-dev:
	docker compose up -d

logs-web:
	docker logs -f store-web

logs-db:
	docker logs -f store-db

check-db:
	docker exec -it store-db psql -U store -d store-db

clear-db:
	docker rm -f -v store-db

migrate-db:
	migrate -database "postgres://store:password@localhost:5432/store-db?sslmode=disable" -path db/migrations up
	
migrate-db-down:
	migrate -database "postgres://store:password@localhost:5432/store-db?sslmode=disable" -path db/migrations down -all
	
build-prod-linux:
	GOOS=linux GOARCH=amd64 go build -o build/eniqilo-store

build-prod-win:
	GOOS=windows GOARCH=amd64 go build -o build/eniqilo-store.exe

build-prod-mac:
	GOOS=darwin GOARCH=amd64 go build -o build/eniqilo-store

scp:
	scp -i w1key build/eniqilo-store ubuntu@xx.xx.xx.xx:~
