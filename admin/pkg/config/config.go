package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpPort              string `env-default:"8080" env:"PORT"`
	CloudProject          string `env-default:"events-consumer"`
	FirestoreCollectionId string `env-default:"jobs"`
	BasicAuthUser         string `env-default:"admin" env:"BASIC_AUTH_USER"`
	BasicAuthPassword     string `env-default:"admin" env:"BASIC_AUTH_PASSWORD"`
	TemplatePath          string `env-default:"./" env:"TEMPLATE_PATH"`
}

var C Config

func init() {
	err := cleanenv.ReadEnv(&C)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Config %+v", C)
}
