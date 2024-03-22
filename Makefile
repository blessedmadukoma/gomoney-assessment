include .env

DB_DRIVER := $(DB_DRIVER)
DB_USER := $(DB_USER)
DB_PASSWORD := $(DB_PASSWORD)
DB_HOST := $(DB_HOST)
DB_PORT := $(DB_PORT)
DB_NAME := $(DB_NAME)

# db/createmigration:
# 	migrate create -ext sql -dir db/migrations -seq $(name)

# db/migrateup:
# 	# migrate up "$(DB_URL)"
# 	migrate -path db/migrations -database "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

# db/migratedown:
# 	# migrate down
# 	migrate -path db/migrations -database "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

dc_up:
# create mongodb and redis containers
	docker-compose up -d

dc_down:
# delete mongodb and redis containers
	docker-compose down -v

# db_up:
# 	# create database
# 	docker exec -it fintrax_db createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)
# 	docker exec -it fintrax_db_live createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

# db_down:
# 	# drop database
# 	docker exec -it fintrax_db dropdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)
# 	docker exec -it fintrax_db_live dropdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

server:
		go run main.go
		# air

air: docs
	air

test:
	go test -v -cover ./...

.PHONY:	db/createmigration db/migrateup db/migratedown pg_up pg_down db_up db_down sqlc start test