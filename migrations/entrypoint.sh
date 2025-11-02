#!/bin/bash

echo "Current working directory: $(pwd)"
echo "Database connstring: postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

for i in {1..5}; do
    migrate -path . -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up
    if [ $? -eq 0 ]; then
        echo "Migrations applied successfully!"
        exit 0
    fi
    echo "Migration attempt $i failed. Retrying in 3 seconds..."
    sleep 3
done

migrate -path . -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" version
# migrate -path . -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up
