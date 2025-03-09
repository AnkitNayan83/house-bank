postgresrun:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:17-alpine

postgresstart:
	docker start postgres17

postgresstop:
	docker stop postgres17

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root house_bank

dropdb:
	docker exec -it postgres17 dropdb house_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/house_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/house_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate

.PHONY: postgresrun postgresstart postgresstop createdb dropdb migrateup migratedown sqlc