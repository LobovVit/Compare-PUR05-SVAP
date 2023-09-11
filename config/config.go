package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Мode       string `yaml:"mode"  env-required:"true"`
	Mastertype string `yaml:"mastertype" env-required:"true"`
	Slavetype  string `yaml:"slavetype" env-required:"true"`
	Masterdsn  string `yaml:"masterdsn"`
	Slavedsn   string `yaml:"slavedsn"`
	Attrs      string `yaml:"attrs"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("read config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); /*cleanenv.ReadEnv(instance)*/ err != nil {
			helptext := "Test text for read config"
			help, _ := cleanenv.GetDescription(instance, &helptext)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
