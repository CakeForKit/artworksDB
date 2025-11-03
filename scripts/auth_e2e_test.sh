#!/bin/bash


# set -a  # Automatically export all variables
# source ./configs/email.env
# set +a
# go test -v github.com/CakeForKit/artworksDB.git/internal/tests/e2e/e2e_auth

# exit 0

docker compose -f ./deployment/docker-compose.ci.yml \
      --env-file ./configs/test_db.env \
      exec \
      -e ALLURE_OUTPUT_PATH=/app \
      -e ALLURE_OUTPUT_DIR=allure-results \
      -e APP_EMAIL_PASSWORD=t5rEDO1haK1YxUdZUrjW \
      -e APP_EMAIL=ktestapp@mail.ru \
      -e TEST_USER_EMAIL=tmpforread@mail.ru \
      -e TEST_USER_EMAIL_PASSWORD=KouQGXYtAXiO73mBcdk6 \
      -T test-runner \
      go test -v ./internal/tests/e2e/e2e_auth