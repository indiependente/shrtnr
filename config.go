package main

import "github.com/caarlos0/env/v7"

type config struct {
	MongoDBPort       string `env:"MONGODB_PORT" envDefault:"27017"`
	MongoDBUser       string `env:"MONGODB_USER" envDefault:"frank"`
	MongoDBPassword   string `env:"MONGODB_PASSWORD" envDefault:"password"`
	MongoDBHost       string `env:"MONGODB_HOST" envDefault:"localhost"`
	MongoDBName       string `env:"MONGODB_DBNAME" envDefault:"shrtnr"`
	MongoDBCollection string `env:"MONGODB_COLLECTION" envDefault:"urls"`
	Port              int    `env:"PORT" envDefault:"7000"`
	SlugLen           int    `env:"SLUG_LEN" envDefault:"5"`
}

func parseConfig() (*config, error) {
	var c config

	return &c, env.Parse(&c)
}
