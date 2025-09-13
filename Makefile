.PHONY: test-allure report open clean

ALLURE_RESULTS_DIR := $(shell pwd)/allure-results
ALLURE_REPORT_DIR := $(shell pwd)/allure-report
ALLURE_HISTORY_DIR := $(shell pwd)/allure-history

allure: report-allure open clean

test-allure:
	ALLURE_LAUNCH_START=$(date +%s000) \
	ALLURE_LAUNCH_END=$(date +%s000) \
	ALLURE_OUTPUT_PATH=$(shell pwd) \
	ALLURE_LAUNCH_NAME="build-$(shell date +%Y%m%d-%H%M%S)" \
	go test -v git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/... 


allure-copy-history:
	@echo "Copy from history to results..."
	mkdir -p $(ALLURE_RESULTS_DIR)/history
	@ if [ -d "$(ALLURE_HISTORY_DIR)" ] && [ -n "$$(ls -A $(ALLURE_HISTORY_DIR) 2>/dev/null)" ]; then \
		cp -r "$(ALLURE_HISTORY_DIR)"/history/* "$(ALLURE_RESULTS_DIR)/history/"; \
		echo "History copied successfully"; \
	else \
		echo "No history found in $(ALLURE_HISTORY_DIR), starting fresh"; \
	fi

report-allure:  test-allure allure-copy-history
	allure generate $(shell pwd)/allure-results -o $(ALLURE_REPORT_DIR) --clean

	@echo "Saving history for next runs..."
	mkdir -p $(ALLURE_HISTORY_DIR)
	@ if [ -d "$(ALLURE_REPORT_DIR)/history" ]; then \
		cp -r "$(ALLURE_REPORT_DIR)/history" "$(ALLURE_HISTORY_DIR)/"; \
		echo "History saved to $(ALLURE_HISTORY_DIR)"; \
	else \
		echo "Warning: No history directory found in report"; \
	fi
	
	mkdir -p $(ALLURE_HISTORY_DIR)
	cp -r "$(ALLURE_REPORT_DIR)/history" "$(ALLURE_HISTORY_DIR)/" 2>/dev/null || true

open:
	allure open ./allure-report
	

clean:
	rm -rf $(ALLURE_RESULTS_DIR) $(ALLURE_REPORT_DIR) 



