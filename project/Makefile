DB_URL=postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable&timezone=UTC&connect_timeout=5

postgres:
	docker run --name postgres_simple_bank -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:14.0

create_db:
	docker exec -it postgres_simple_bank createdb --username=postgres --owner=postgres simple_bank

drop_db:
	docker exec -it postgres_simple_bank dropdb simple_bank

migrate_up:
	migrate -path ../db/migration -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path ../db/migration -database "$(DB_URL)" -verbose down

sqlc:
	cd ../ && sqlc generate

test:
	go test -v -cover ../...

.PHONY: postgres createdb dropdb migrate_up migrate_down sqlc
