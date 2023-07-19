include .env

## Delete consumer
consumer-rm:
	docker rm -f report-service
	docker rmi report-image

## Build and run consumer
consumer-up:
	docker build -t report-image .
	docker run -dp 8081:8081\
		--network order-and-pay-app --network-alias report\
		-e POSTGRES_HOST=${POSTGRES_HOST}\
		-e POSTGRES_USER=${POSTGRES_USER}\
		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD}\
		-e POSTGRES_DB=${POSTGRES_DB}\
		-e POSTGRES_PORT=${POSTGRES_PORT}\
		--name report-service\
		report-image

## Run migrations UP
migrate-up:
	docker compose run --rm migrate up

## Run migrations DOWN
migrate-down:
	docker compose run --rm migrate down 1

## Create a DB migration files e.g `make migrate-create name=migration-name`
migrate-create:
	docker compose run --rm migrate create -ext sql -dir /migrations -seq $(name)

## Enter to database console
shell-db:
	docker compose exec postgre psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}