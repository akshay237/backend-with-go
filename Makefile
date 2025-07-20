network:
	docker network create bank-network

postgres:
	docker run --name postgres-17 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it postgres-17 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -i postgres-17 dropdb simple_bank

migrateup:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	set -e; \
	go run main.go

mock:
	mockgen -package mockdb -destination database/mock/store.go github.com/akshay237/backend-with-go/database/sqlc Store

shutdown:
	touch ./stopfile

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock shutdown proto