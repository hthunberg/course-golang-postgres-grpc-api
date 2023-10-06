DOCKER_COMPOSE_FILE ?= ./docker/docker-compose.yaml

#****************
# DB migrations
#****************

.PHONY: migrate-db-up
migrate-db-up:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm db-migrate up

.PHONY: migrate-db-down
migrate-db-down:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm db-migrate down 1

.PHONY: migrate-db-create-next-migration-path
migrate-db-create-next-migration-path:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm db-migrate create -ext sql -dir ./migrations bank-db-migration

#****************
# compose landscape
#****************

.PHONY: build up start down destroy stop restart logs ps db-shell
build:
	docker compose -f ${DOCKER_COMPOSE_FILE} build $(c)
up:
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d $(c)
start:
	docker compose -f ${DOCKER_COMPOSE_FILE} start $(c)
down:
	docker compose -f ${DOCKER_COMPOSE_FILE} down $(c)
destroy:
	docker compose -f ${DOCKER_COMPOSE_FILE} down -v $(c)
stop:
	docker compose -f ${DOCKER_COMPOSE_FILE} stop $(c)
restart:
	docker compose -f ${DOCKER_COMPOSE_FILE} stop $(c)
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d $(c)
logs:
	docker compose -f ${DOCKER_COMPOSE_FILE} logs --tail=100 -f $(c)
ps:
	docker compose -f ${DOCKER_COMPOSE_FILE} ps
db-shell:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec db psql -U postgres -d postgres