package backend

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ImpactUsername string
	ImpactPassword string
	CallbackUsername string
	CallbackPassword string
	CallbackHost string
}

func ReadConfig() *Config {
	fp, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Reading config file", err)
	}
	cfg := Config{}
	err = json.Unmarshal(fp, &cfg)
	if err != nil {
		log.Fatalln("Parsing config file", err)
	}
	return &cfg

}