package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	Database string `default:"postgres"`
}

type Client struct {
	PgxConn *pgx.Conn
}

func NewEnvConfig(service_name string) (Config, error) {
	var c Config

	return c, envconfig.Process(service_name, &c)
}

func NewClient(conf Config) (Client, error) {
	dsnPgx := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

	conn, err := pgx.Connect(context.Background(), dsnPgx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		// return nil, err
		// os.Exit(1)
	}

	return Client{PgxConn: conn}, nil
}
