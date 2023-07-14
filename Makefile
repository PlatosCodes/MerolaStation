postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=bluecomet -d postgres:15-alpine
	
createdb: 
	docker exec -it postgres15 createdb --username=root --owner=root merola_station

dropdb:
	docker exec -it postgres15 dropdb --username=root merola_station

citext:
	docker exec -it postgres15 psql --username=root merola_station -c "CREATE EXTENSION IF NOT EXISTS citext;"

migrateup:
	migrate -path db/migration -database "postgresql://root:bluecomet@localhost:5432/merola_station?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://root:bluecomet@localhost:5432/merola_station?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb citext migrateup migratedown sqlc test