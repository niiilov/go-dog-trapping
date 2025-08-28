package config

type Config struct {
	*Postgres
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

func NewConfig() *Config {

	pg := Postgres{}
	pg.Host = "82.202.169.245"
	pg.Port = 5432
	pg.User = "admin"
	pg.Password = "1234"
	pg.Database = "sobaki"
	pg.MaxConn = 10
	pg.MinConn = 5

	return &Config{
		Postgres: &pg,
	}
}
