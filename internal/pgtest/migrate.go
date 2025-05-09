package pgtest

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// func MigrateUp(ctx context.Context, db *cnfg.PostgresCredentials, migrationDir string) error {
func MigrateUp(ctx context.Context, migrationDir string, pgCreds *cnfg.PostgresCredentials) error {
	// wd, err := os.Getwd() // Получает директорию, из которой запущен `go run`
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Working directory:", wd)
	sourceUrl := fmt.Sprintf("file://%s", migrationDir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	m, err := migrate.New(sourceUrl, dbUrl)
	fmt.Printf("Migrations: sourceUrl=%s, dbUrl=%s\n", sourceUrl, dbUrl)
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
func MigrateDown(ctx context.Context, migrationDir string, pgCreds *cnfg.PostgresCredentials) error {
	sourceUrl := fmt.Sprintf("file://%s", migrationDir)
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
