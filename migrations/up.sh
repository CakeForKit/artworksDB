#!/bin/bash

migrate -path ./internal/migrations -database "postgresql://puser:ppassword@localhost:5432/artworks?sslmode=disable" -verbose up