
ALLURE_OUTPUT_PATH := $(shell pwd)
ALLURE_RESULTS_DIR := $(shell pwd)/allure-results
ALLURE_REPORT_DIR := $(shell pwd)/allure-report
export ALLURE_RESULTS_DIR
export ALLURE_OUTPUT_PATH

SCRIPTS := ./scripts
DC_CI := ./deployment/docker-compose.ci.yml
TEST_DB_ENV := ./configs/test_db.env

.PHONY: allure
allure: unit_test integration_test report_allure open_allure

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
	docker compose -v -f $(DC_CI) --env-file $(TEST_DB_ENV) build --no-cache --progress=plain test-runner
	docker compose -v -f $(DC_CI) --env-file $(TEST_DB_ENV) up  test-runner

.PHONY: down_app
down_app:
	docker compose -f $(DC_CI) down -v test-runner


.PHONY: run_serv
run_serv:
	docker compose -f $(DC_CI) --env-file $(TEST_DB_ENV) up -d postgres migrator redis_artworks

.PHONY: down_serv
down_serv:
	docker compose -f $(DC_CI) down -v postgres migrator redis_artworks

.PHONY: build
build:
	docker compose -f $(DC_CI) --env-file $(TEST_DB_ENV) build


# .PHONY: test-allure report open clean

# ALLURE_RESULTS_DIR := $(shell pwd)/allure-results
# ALLURE_REPORT_DIR := $(shell pwd)/allure-report
# ALLURE_HISTORY_DIR := $(shell pwd)/allure-history

# allure: report-allure open clean

# test-allure:
# 	ALLURE_LAUNCH_START=$(date +%s000) \
# 	ALLURE_LAUNCH_END=$(date +%s000) \
# 	ALLURE_OUTPUT_PATH=$(shell pwd) \
# 	ALLURE_LAUNCH_NAME="build-$(shell date +%Y%m%d-%H%M%S)" \
# 	go test -v -shuffle=on git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/...; \
# 	SERVICES_EXIT_CODE=$$?; \
# 	if [ $$SERVICES_EXIT_CODE -eq 0 ]; then \
# 		docker compose -f ./deployment/docker-compose.ci.yml \
# 			--env-file ./configs/test_db.env \
# 			exec \
# 			-e ALLURE_LAUNCH_START=$(date +%s000) \
# 			-e ALLURE_LAUNCH_END=$(date +%s000) \
# 			-e ALLURE_OUTPUT_PATH=/app \
# 			-e ALLURE_LAUNCH_NAME="build-$(date +%Y%m%d-%H%M%S)" \
# 			test-runner \
# 			go test -v -shuffle=on git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/integration; \
# 			INTEGRATION_EXIT_CODE=$$?; \
# 	fi; \
# 	exit 0

# allure-copy-history:
# 	@echo "Copy from history to results..."
# 	mkdir -p $(ALLURE_RESULTS_DIR)/history
# 	@ if [ -d "$(ALLURE_HISTORY_DIR)" ] && [ -n "$$(ls -A $(ALLURE_HISTORY_DIR) 2>/dev/null)" ]; then \
# 		cp -r "$(ALLURE_HISTORY_DIR)"/history/* "$(ALLURE_RESULTS_DIR)/history/"; \
# 		echo "History copied successfully"; \
# 	else \
# 		echo "No history found in $(ALLURE_HISTORY_DIR), starting fresh"; \
# 	fi

# report-allure: test-allure allure-copy-history
# 	allure generate $(shell pwd)/allure-results -o $(ALLURE_REPORT_DIR) --clean

# 	@echo "Saving history for next runs..."
# 	mkdir -p $(ALLURE_HISTORY_DIR)
# 	@ if [ -d "$(ALLURE_REPORT_DIR)/history" ]; then \
# 		cp -r "$(ALLURE_REPORT_DIR)/history" "$(ALLURE_HISTORY_DIR)/"; \
# 		echo "History saved to $(ALLURE_HISTORY_DIR)"; \
# 	else \
# 		echo "Warning: No history directory found in report"; \
# 	fi
	
# 	mkdir -p $(ALLURE_HISTORY_DIR)
# 	cp -r "$(ALLURE_REPORT_DIR)/history" "$(ALLURE_HISTORY_DIR)/" 2>/dev/null || true

# open:
# 	allure open ./allure-report
	

# run_serv:
# 	docker compose -f ./deployment/docker-compose.ci.yml --env-file ./configs/test_db.env up -d postgres migrator test-runner

# build:
# 	docker compose -f ./deployment/docker-compose.ci.yml --env-file ./configs/test_db.env build

# down_serv:
# 	docker compose -f ./deployment/docker-compose.ci.yml down -v

# clean:
# 	rm -rf $(ALLURE_RESULTS_DIR) $(ALLURE_REPORT_DIR) 





# run:
# 	docker compose -f ./deployment/docker-compose.ci.yml --env-file ./configs/test_db.env up --build  test-runner
# # 	docker compose -f ./deployment/docker-compose.ci.yml --env-file ./configs/test_db.env run -T --rm test-runner go test -v git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/integration 

# up:
# 	docker compose -f ./deployment/docker-compose.ci.yml --env-file ./configs/test_db.env up --build -d postgres migrator

# clear:
# 	docker compose -f ./deployment/docker-compose.ci.yml down -v