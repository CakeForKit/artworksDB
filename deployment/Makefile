DOCKER_COMPOSE = docker-compose
PROD_FILE = docker-compose.yaml
DEV_FILE = docker-compose.dev.yaml

prod:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) build --no-cache app

prod_up:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) up app

dev:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) build --no-cache app
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up app

dev_up:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up app

clickhouse:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up --build clickhouse_artworks clickhouse_migrator_artworks

down:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) down -v clickhouse_artworks clickhouse_migrator_artworks

all:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up --build -d

clear_templ:
	rm -f ./internal/frontend/components/*_templ.go