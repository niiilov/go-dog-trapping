package config

import (
	"os"
	"strconv"
)

type Config struct {
	*Postgres
	*AwsCreds
}

type Postgres struct {
	Host     string `env:"PG_HOST"`
	Port     int    `env:"PG_PORT"`
	Database string `env:"PG_DATABASE"`
	User     string `env:"PG_USER"`
	Password string `env:"PG_PASS"`
	MaxConn  int32  `env:"PG_MAXCONN"`
	MinConn  int32  `env:"PG_MINCONN"`
}

type AwsCreds struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BaseEndpoint    string
}

func NewConfig() *Config {

	pg := Postgres{}
	pg.Host = os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	pg.Port, _ = strconv.Atoi(port)

	pg.User = os.Getenv("PG_USER")
	pg.Password = os.Getenv("PG_PASS")
	pg.Database = os.Getenv("PG_DATABASE")
	pg.MaxConn = 10
	pg.MinConn = 5

	aws := AwsCreds{}
	aws.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	aws.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	aws.BaseEndpoint = os.Getenv("AWS_BASE_ENDPOINT")
	aws.Region = os.Getenv("AWS_REGION")

	return &Config{
		Postgres: &pg,
		AwsCreds: &aws,
	}
}
