package pgtest

import (
	"context"
	"fmt"
	"sync"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pgContainer  testcontainers.Container
	pgCreds      cnfg.PostgresCredentials
	pgTestConfig cnfg.PostgresTestConfig
	pgOnce       sync.Once
	pgSetupErr   error
)

func GetTestPostgres(ctx context.Context) (testcontainers.Container, cnfg.PostgresCredentials, error) {
	pgOnce.Do(func() {
		pgTestConfig = *cnfg.GetPgTestConfig()
		pgContainer, pgCreds, pgSetupErr = NewTestPostgres(ctx, &pgTestConfig)
	})
	return pgContainer, pgCreds, pgSetupErr
}

func NewTestPostgres(ctx context.Context, config *cnfg.PostgresTestConfig) (testcontainers.Container, cnfg.PostgresCredentials, error) {
	strPort := fmt.Sprintf("%d/tcp", config.Port)
	// strPort := "5432/tcp"
	fmt.Printf("NewTestPostgres: 0\n")
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		ExposedPorts: []string{strPort},
		Env: map[string]string{
			"POSTGRES_USER":     config.Username,
			"POSTGRES_PASSWORD": config.Password,
			"POSTGRES_DB":       config.DbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port(strPort)),
		),
		AutoRemove: true,
	}
	fmt.Printf("NewTestPostgres: %+v\n", config)
	cnt, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	fmt.Printf("NewTestPostgres: genered\n")
	if err != nil {
		return nil, cnfg.PostgresCredentials{}, fmt.Errorf("NewTestPostgres: %w", err)
	}
	fmt.Printf("NewTestPostgres: 1\n")
	host, err := cnt.Host(ctx)
	fmt.Printf("NewTestPostgres: 1.1\n")
	if err != nil {
		return nil, cnfg.PostgresCredentials{}, fmt.Errorf("NewTestPostgres: %w", err)
	}
	fmt.Printf("NewTestPostgres: host - %s\n", host)
	port, err := cnt.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, cnfg.PostgresCredentials{}, fmt.Errorf("NewTestPostgres: %w", err)
	}
	creds := cnfg.PostgresCredentials{
		Host:     host,
		DbName:   config.DbName,
		Port:     port.Int(),
		Username: config.Username,
		Password: config.Password,
	}
	// creds := NewPostgresCredentials(config.User, config.Password, config.Database, host, uint16(port.Int()))
	return cnt, creds, nil
}
