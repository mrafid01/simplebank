DB_URL=postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable

postgres:
	docker run --name postgres17 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres17 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

dbdocs:
	dbdocs build doc/db.dbml

dbschema:
	dbml2sql doc/db.dbml --postgres -o doc/schema.sql

db_generate_from_database:
	db2dbml postgres "$(DB_URL)" -o doc/schema_from_database.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/mrafid01/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 9090 --package pb --service SimpleBank -r repl

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 dbdocs dbschema db_generate_from_database sqlc server mock proto evans