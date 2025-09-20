package config

import (
	"os"
)

type Config struct {
	*Postgres
	*AwsCreds
}

type Postgres struct {
	Host     string `env:"PG_HOST"`
	Port     int    `env:"PG_PORT"`
	Database string `env:"PGPG_DATABASE_DATA"`
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
	pg.Host = "82.202.169.245"
	pg.Port = 5432
	pg.User = "admin"
	pg.Password = "1234"
	pg.Database = "sobaki"
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
