package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HttpPort              string `default:"8080",envconfig:"port"`
	CloudProject          string `default:"events-consumer"`
	FirestoreCollectionId string `default:"jobs"`
	BasicAuthUser         string `default:"admin",envconfig:"basic_auth_user"`
	BasicAuthPassword     string `default:"admin",envconfig:"basic_auth_password"`
}

var C Config

func init() {
	err := envconfig.Process("", &C)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Config %+v", C)
}
