package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string `yaml:"env"  env-default:"local"`
	TgBotToken      string `yaml:"tgbottoken" env-required:"true"`
	YandexDiskToken string `yaml:"yandexdisktoken" env-required:"true"`
}

func MustLoad() *Config {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(res); os.IsNotExist(err) {
		panic("config file does not exist: " + res)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(res, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg
}
