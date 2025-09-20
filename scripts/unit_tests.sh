#!/bin/bash

ALLURE_LAUNCH_START=$(date +%s000) \
ALLURE_LAUNCH_END=$(date +%s000) \
ALLURE_LAUNCH_NAME="unit-test-$(shell date +%Y%m%d-%H%M%S)" \
go test -shuffle=on git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/userservice;
exit 0
