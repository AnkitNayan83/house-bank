DB_URL=postgresql://root:postgres@localhost:5432/house_bank?sslmode=disable

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
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateupone:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migratedownone:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

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

schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

dbdos:
	dbdocs build ./doc/db.dbml

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto
evans:
	evans --host localhost --port 8000 --reflection --repl


.PHONY: postgresconsole image postgresrun postgresstart postgresstop createdb dropdb newmigration migrateup migrateupone migratedown migratedownone sqlc test server mock proto evans