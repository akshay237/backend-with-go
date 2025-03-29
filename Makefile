postgres:
	docker run --name postgres-17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it postgres-17 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -i postgres-17 dropdb simple_bank

migrateup:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/akshay237/backend-with-go/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock