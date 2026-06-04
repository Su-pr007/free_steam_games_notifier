-include .env
export

build:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.build.yaml build

up:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml up -d

down:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml down

release:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml down
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml up -d

ps:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml ps -a

logs:
	$(DOCKER_COMPOSE_COMMAND) -f docker-compose.prod.yaml logs

push:
	docker push ghcr.io/su-pr007/freesteamgamesparser-go:latest

