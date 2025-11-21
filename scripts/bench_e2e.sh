#!/bin/bash

set -a  # Automatically export all variables
source ./configs/email.env
set +a

start=$(date -u +%Y-%m-%dT%H:%M:%SZ)   

ALLURE_LAUNCH_START=$(date +%s000) \
ALLURE_LAUNCH_END=$(date +%s000) \
ALLURE_LAUNCH_NAME="e2e-test-$(shell date +%Y%m%d-%H%M%S)" \
go test -shuffle=on github.com/CakeForKit/artworksDB.git/internal/tests/e2e/e2e_api;

end=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go run ./cmd/export_prom/main.go -start="$start" -end="$end" -container="app_artworks_dev"

exit 0