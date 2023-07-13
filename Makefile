postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=bluecomet -d postgres:15-alpine
	
createdb: 
	docker exec -it postgres15 createdb --username=root --owner=root merola_station

dropdb:
	docker exec -it postgres15 dropdb --username=root --owner=root merola_station

migrateup:
	migrate -path db/migration -database "postgresql://root:bluecomet@localhost:5432/merola_station?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://root:bluecomet@localhost:5432/merola_station?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc