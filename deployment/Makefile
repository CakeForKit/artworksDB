DOCKER_COMPOSE = docker-compose
PROD_FILE = docker-compose.yaml
DEV_FILE = docker-compose.dev.yaml

prod:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) up --build app

prod_up:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) up app

dev:
	$(DOCKER_COMPOSE) -f $(PROD_FILE) -f $(DEV_FILE) up --build app
