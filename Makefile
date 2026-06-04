-include .env
export

up:
	$(DOCKER_COMPOSE_COMMAND) up -d

build:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.build.yaml build

up-prod:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml up -d

release-prod:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml down
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml up -d

ps-prod:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml ps -a

logs-prod:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml logs

push:
	docker push ghcr.io/su-pr007/freesteamgamesparser-go:latest

