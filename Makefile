postgres:
	docker container run -dt --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres postgres:12-alpine

startdb:
	docker container start postgres

stopdb:
	docker container stop postgres

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simplebank_db

createdb_test:
	docker exec -it postgres createdb --username=postgres --owner=postgres simplebank_test_db

dropdb:
	docker exec -it postgres dropdb --username=postgres simplebank_db

dropdb_test:
	docker exec -it postgres dropdb --username=postgres simplebank_test_db

migrateup:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/simplebank_db?sslmode=disable" -verbose up

migrateup_test:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/simplebank_test_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/simplebank_db?sslmode=disable" -verbose down

migratedown_test:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/simplebank_test_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres startdb stopdb createdb createdb_test dropdb dropdb_test migrateup migrateup_test migratedown migratedown_test sqlc test
