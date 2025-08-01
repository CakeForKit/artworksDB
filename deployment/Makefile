DOCKER_COMPOSE = docker-compose
PROD_FILE = docker-compose.yaml
DEV_FILE = docker-compose.dev.yaml
SERVS = postgres migrator redis_artworks # clickhouse_artworks clickhouse_migrator_artworks

# prod:
# 	$(DOCKER_COMPOSE) -f $(PROD_FILE) build --no-cache app

# prod_up:
# 	$(DOCKER_COMPOSE) -f $(PROD_FILE) up app

prep:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) build app $(SERVS)
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up $(SERVS) -d

up:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up app

build_up:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) build --no-cache app
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up app

clickhouse:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up --build clickhouse_artworks clickhouse_migrator_artworks

down:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) down -v clickhouse_artworks clickhouse_migrator_artworks

all:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up --build 

clear_templ:
	rm -f ./internal/frontend/components/*_templ.go