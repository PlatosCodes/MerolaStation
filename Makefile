DB_URL=postgresql://root:bluecomet@localhost:5432/merolastation?sslmode=disable

network:
	docker network create merolastation-network
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=bluecomet -d postgres
	
createdb: 
	docker exec -it postgres createdb --username=root merolastation

dropdb:
	docker exec -it postgres dropdb --username=root merolastation

migrateup:
	migrate -path db/migration -database "${DB_URL}" -verbose up

migrateup1:
	migrate -path db/migration -database "${DB_URL}" -verbose up 1

migratedown: 
	migrate -path db/migration -database "${DB_URL}" -verbose down

migratedown1: 
	migrate -path db/migration -database "${DB_URL}" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

resetdb:
	(docker exec -it postgres dropdb --username=root merolastation || true) \
	&& docker exec -it postgres createdb --username=root merolastation \
	&& migrate -path db/migration -database "${DB_URL}" -verbose up

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/PlatosCodes/MerolaStation/db/sqlc Store \
	&& mockgen -package mockmailer -destination mailer/mock/mailer.go github.com/PlatosCodes/MerolaStation/mailer IMailer
	
.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 new_migration db_docs db_schema sqlc test server resetdb mock