
ALLURE_OUTPUT_PATH := $(shell pwd)
ALLURE_RESULTS_DIR := $(shell pwd)/allure-results
ALLURE_REPORT_DIR := $(shell pwd)/allure-report
export ALLURE_RESULTS_DIR
export ALLURE_OUTPUT_PATH

SCRIPTS := ./scripts
DC_CI := ./deployment/docker-compose.ci.yml
TEST_DB_ENV := ./configs/test_db.env
DC_DEV := ./deployment/docker-compose.dev.yaml
DB_ENV := ./configs/db.env

.PHONY: allure
allure: unit_test integration_test e2e_test report_allure open_allure

.PHONY: unit_test
unit_test : clear_allure
	$(SCRIPTS)/unit_tests.sh

.PHONY: integration_test
integration_test: clear_allure
	$(SCRIPTS)/integration_tests.sh

.PHONY: e2e_test
e2e_test: clear_allure
	$(SCRIPTS)/e2e_test.sh

.PHONY: report_allure
report_allure:
	mkdir -p $(ALLURE_REPORT_DIR)/history
	cp -r $(ALLURE_REPORT_DIR)/history $(ALLURE_RESULTS_DIR)
	allure generate $(ALLURE_RESULTS_DIR) -o $(ALLURE_REPORT_DIR) --clean

.PHONY: clear_allure
clear_allure:
	rm -rf $(ALLURE_RESULTS_DIR)

.PHONY: open_allure
open_allure:
	allure open $(ALLURE_REPORT_DIR)
	


.PHONY: run_app
run_app:
	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) up --build app

.PHONY: down_app
down_app:
	docker compose -f $(DC_DEV) down -v app

.PHONY: run_serv
run_serv:
	docker compose -f $(DC_DEV) --env-file $(DB_ENV) up -d postgres migrator redis_artworks

.PHONY: down_serv
down_serv:
	docker compose -f $(DC_DEV) --env-file $(DB_ENV) down -v postgres migrator redis_artworks

.PHONY: stop_serv
stop_serv:
	docker compose -f $(DC_DEV) --env-file $(DB_ENV) stop  postgres migrator redis_artworks


.PHONY: run_test_app
run_test_app:
# --no-cache
	docker compose -v -f $(DC_CI) --env-file $(TEST_DB_ENV) build --progress=plain test-runner
	docker compose -v -f $(DC_CI) --env-file $(TEST_DB_ENV) up  test-runner

.PHONY: down_test_app
down_test_app:
	docker compose -f $(DC_CI) down -v test-runner

.PHONY: run_test_serv
run_test_serv:
	docker compose -f $(DC_CI) --env-file $(TEST_DB_ENV) up -d postgres migrator redis_artworks

.PHONY: down_test_serv
down_test_serv:
	docker compose -f $(DC_CI) down -v postgres migrator redis_artworks

.PHONY: build
build:
	docker compose -f $(DC_CI) --env-file $(TEST_DB_ENV) build


.PHONY: clear_docker
clear_docker:
# Остановите все контейнеры
	docker-compose -f ./deployment/docker-compose.ci.yml down
# Удалите старые образы
	docker rmi deployment-test-runner
# Очистите builder кеш
	docker builder prune -f
# Удалите все старые версии
	docker rmi deployment-test-runner:latest
# Полная очистка
	docker system prune -a -f


