include .env

DB_DRIVER := $(DB_DRIVER)
DB_USER := $(DB_USER)
DB_PASSWORD := $(DB_PASSWORD)
DB_HOST := $(DB_HOST)
DB_PORT := $(DB_PORT)
DB_NAME := $(DB_NAME)

dc_up:
# create mongodb and redis containers
	docker-compose up -d

dc_down:
# delete mongodb and redis containers
	docker-compose down -v

seed:
	go run main.go -seed

server:
		go run main.go
		# air

air: docs
	air

test:
	go test -v -cover ./...

docker-build-and-push:
	docker build -t gomoney-premier-league .
	
	docker images
	
	docker tag $(shell docker images -q gomoney-premier-league) blessedmadukoma/gomoney-assessment

	docker login
    
	docker push blessedmadukoma/gomoney-assessment:latest

.PHONY:dc_up dc_down sqlc start test server seed docker-build-and-push