postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=bluecomet -d postgres:15-alpine
	
createdb: 
	docker exec -it postgres15 createdb --username=root --owner=root merola_station

dropdb:
	docker exec -it postgres15 dropdb --username=root merola_station

citext:
	docker exec -it postgres15 psql --username=root merola_station -c "CREATE EXTENSION IF NOT EXISTS citext;"

migrateup:
	migrate -path db/migration -database "postgresql://root:bluecomet@localhost:5432/merola_station?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:bluecomet@localhost:5432/merola_station?sslmode=disable" -verbose up 1

migratedown: 
	migrate -path db/migration -database "postgresql://root:bluecomet@localhost:5432/merola_station?sslmode=disable" -verbose down

migratedown1: 
	migrate -path db/migration -database "postgresql://root:bluecomet@localhost:5432/merola_station?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/PlatosCodes/MerolaStation/db/sqlc Store
	
.PHONY: postgres createdb dropdb citext migrateup migrateup1 migratedown migratedown1 sqlc test server mock