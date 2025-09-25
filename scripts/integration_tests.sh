#!/bin/bash

ALLURE_LAUNCH_START=$(date +%s000) \
ALLURE_LAUNCH_END=$(date +%s000) \
ALLURE_LAUNCH_NAME="unit-test-$(shell date +%Y%m%d-%H%M%S)" \
go test -v -run "TestEventRepSuite/EventRepSuite/" git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/integration/...;
exit 0

ALLURE_LAUNCH_START=$(date +%s000) \
ALLURE_LAUNCH_END=$(date +%s000) \
ALLURE_LAUNCH_NAME="unit-test-$(shell date +%Y%m%d-%H%M%S)" \
go test -v  git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/integration/...;
# -shuffle=on

# docker compose -f ./deployment/docker-compose.ci.yml \
# 			--env-file ./configs/test_db.env \
# 			exec \
#             -e ALLURE_OUTPUT_PATH=/app \
# 			-e ALLURE_LAUNCH_START=$(date +%s000) \
# 			-e ALLURE_LAUNCH_END=$(date +%s000) \
# 			-e ALLURE_LAUNCH_NAME="build-$(date +%Y%m%d-%H%M%S)" \
# 			test-runner \
# 			go test -shuffle=on git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/integration; 

exit 0
