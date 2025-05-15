migrateup:
	migrate -path ./internal/migrations -database "postgresql://puser:ppassword@postgres_container:5432/artworks?sslmode=disable" -verbose up

migratedown:
	migrate -path ./internal/migrations -database "postgresql://puser:ppassword@postgres_container:5432/artworks?sslmode=disable" -verbose down

.PHONY: migrateup migratedown