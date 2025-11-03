#!/bin/bash


set -a  # Automatically export all variables
source ./configs/email.env
set +a
go test -v github.com/CakeForKit/artworksDB.git/internal/tests/e2e/e2e_auth

exit 0