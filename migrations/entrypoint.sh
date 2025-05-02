#!/bin/bash

echo "Current working directory: $(pwd)"
echo "Directory contents:"
ls -la

migrate -path . -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up