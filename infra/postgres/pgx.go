package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5438"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	Database string `default:"postgres"`
}

type Client struct {
	Pool *pgxpool.Pool
}

func NewEnvConfig(service_name string) (Config, error) {
	var c Config

	return c, envconfig.Process(service_name, &c)
}

func NewClient(conf Config) (Client, error) {
	dsnPgx := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

	poolConf, err := pgxpool.ParseConfig(dsnPgx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		return Client{}, err
	}

	poolConf.MaxConns = 10
	poolConf.MaxConnLifetime = time.Second

	pool, err := pgxpool.NewWithConfig(context.TODO(), poolConf)

	fmt.Println("Connecting to Postgres DB: ", dsnPgx)

	// pool, err := pgx.Connect(context.Background(), dsnPgx)
	// pool, err := pgxpool.New(context.Background(), dsnPgx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(conf)
		// return nil, err
		// os.Exit(1)
	}

	_, acqErr := pool.Acquire(context.TODO())
	if acqErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire connection: %v\n", acqErr)
		panic(conf)
	}

	return Client{Pool: pool}, nil
}
