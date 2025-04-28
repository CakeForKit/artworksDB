package pgtest

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// func MigrateUp(ctx context.Context, db *cnfg.PostgresCredentials, migrationDir string) error {
func MigrateUp(ctx context.Context) error {
	wd, err := os.Getwd() // Получает директорию, из которой запущен `go run`
	if err != nil {
		panic(err)
	}
	fmt.Println("Working directory:", wd)
	sourceUrl := fmt.Sprintf("file://%s", pgTestConfig.MigrationDir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	m, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		return fmt.Errorf("sourceUrl=%s, dbUrl=%s - %w", sourceUrl, dbUrl, err)
	}
	defer m.Close()
	err = m.Up()
	if err != nil {
		return fmt.Errorf("sourceUrl=%s, dbUrl=%s - %w", sourceUrl, dbUrl, err)
	}
	return nil
}

// func MigrateDown(ctx context.Context, db *cnfg.PostgresCredentials, migrationDir string) error {
func MigrateDown(ctx context.Context) error {
	sourceUrl := fmt.Sprintf("file://%s", pgTestConfig.MigrationDir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	m, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		return err
	}
	defer m.Close()
	err = m.Down()
	if err != nil {
		return err
	}
	return nil
}
