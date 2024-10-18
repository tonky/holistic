package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	Database string `default:"postgres"`
}

type PostgresClient struct {
	PgxConn *pgx.Conn
}

func NewPosgresClient(conf PostgresConfig) (*PostgresClient, error) {
	dsnPgx := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

	conn, err := pgx.Connect(context.Background(), dsnPgx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		// return nil, err
		// os.Exit(1)
	}

	return &PostgresClient{PgxConn: conn}, nil
}
