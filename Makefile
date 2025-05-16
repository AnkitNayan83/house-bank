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
	rm -f doc/swagger/*.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=house_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 8000 --reflection --repl

redis:
	docker run --name redis -p 6379:6379 -d redis:8.0.1-alpine


.PHONY: postgresconsole image postgresrun postgresstart postgresstop createdb dropdb newmigration migrateup migrateupone migratedown migratedownone sqlc test server mock proto evans redis