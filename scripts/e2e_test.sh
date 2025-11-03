#!/bin/bash

set -a  # Automatically export all variables
source ./configs/email.env
set +a
ALLURE_LAUNCH_START=$(date +%s000) \
ALLURE_LAUNCH_END=$(date +%s000) \
ALLURE_LAUNCH_NAME="e2e-test-$(shell date +%Y%m%d-%H%M%S)" \
go test -shuffle=on github.com/CakeForKit/artworksDB.git/internal/tests/e2e/e2e_api;

exit 0