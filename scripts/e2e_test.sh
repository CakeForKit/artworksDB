#!/bin/bash

ALLURE_LAUNCH_START=$(date +%s000) \
ALLURE_LAUNCH_END=$(date +%s000) \
ALLURE_LAUNCH_NAME="e2e-test-$(shell date +%Y%m%d-%H%M%S)" \
go test -shuffle=on git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/e2e/e2e_api;

exit 0