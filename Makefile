postgresconsole:
	docker exec -it postgres17 psql -U root -d house_bank
postgresrun:
	docker run --name postgres17 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:17-alpine

postgresstart:
	docker start postgres17

postgresstop:
	docker stop postgres17

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root house_bank

dropdb:
	docker exec -it postgres17 dropdb house_bank

newmigration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/house_bank?sslmode=disable" -verbose up
migrateupone:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/house_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/house_bank?sslmode=disable" -verbose down
migratedownone:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/house_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockDB -destination db/mock/store.go github.com/AnkitNayan83/houseBank/db/sqlc Store

image:
	docker build -t housebank:latest .


.PHONY: postgresconsole image postgresrun postgresstart postgresstop createdb dropdb newmigration migrateup migrateupone migratedown migratedownone sqlc test server mock