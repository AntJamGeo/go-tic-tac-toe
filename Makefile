DEV_COMPOSE_FILE=docker-compose-dev.yml
PROD_COMPOSE_FILE=docker-compose.yml

DEV_IMAGE=go-tic-tac-toe-dev
PROD_IMAGE=go-tic-tac-toe

# Development
.PHONY: build-dev
build-dev:
	docker compose -f $(DEV_COMPOSE_FILE) build

.PHONY: up-dev
up-dev:
	docker compose -f $(DEV_COMPOSE_FILE) up

.PHONY: down-dev
down-dev:
	docker compose -f $(DEV_COMPOSE_FILE) down

.PHONY: clean-dev
clean-dev: down-dev
	docker image rm $(DEV_IMAGE)

# Production
.PHONY: build
build:
	docker compose -f $(PROD_COMPOSE_FILE) build

.PHONY: up
up:
	docker compose -f $(PROD_COMPOSE_FILE) up

.PHONY: down
down:
	docker compose -f $(PROD_COMPOSE_FILE) down

.PHONY: clean
clean: down
	docker image rm $(PROD_IMAGE)

# Clean all
.PHONY: clean-all
clean-all: clean-dev clean