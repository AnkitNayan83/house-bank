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
	migrate -path db/migration -database "postgresql://house_bank_owner:npg_b3kd4pvQoqEy@ep-mute-glitter-a4sca5g7-pooler.us-east-1.aws.neon.tech/house_bank?sslmode=require" -verbose up
migrateupone:
	migrate -path db/migration -database "postgresql://house_bank_owner:npg_b3kd4pvQoqEy@ep-mute-glitter-a4sca5g7-pooler.us-east-1.aws.neon.tech/house_bank?sslmode=require" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://house_bank_owner:npg_b3kd4pvQoqEy@ep-mute-glitter-a4sca5g7-pooler.us-east-1.aws.neon.tech/house_bank?sslmode=require" -verbose down
migratedownone:
	migrate -path db/migration -database "postgresql://house_bank_owner:npg_b3kd4pvQoqEy@ep-mute-glitter-a4sca5g7-pooler.us-east-1.aws.neon.tech/house_bank?sslmode=require" -verbose down 1

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