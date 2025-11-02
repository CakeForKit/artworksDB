#!/bin/bash

ALLURE_LAUNCH_START=$(date +%s000) \
ALLURE_LAUNCH_END=$(date +%s000) \
ALLURE_LAUNCH_NAME="unit-test-$(shell date +%Y%m%d-%H%M%S)" \
go test -shuffle=on github.com/CakeForKit/artworksDB.git/internal/services/...;
exit 0
