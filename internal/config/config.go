package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string      `yaml:"env" env-default:"local"`
	Bot  *BotConfig  `yaml:"bot"`
	Mono *MonoConfig `yaml:"mono"`
	DB   *DBConfig   `yaml:"db"`
}

type BotConfig struct {
	Token      string `yaml:"bot_token"`
	LongPoller int    `yaml:"long_poller"`
}

type MonoConfig struct {
	ApiURL string `yaml:"mono_api_url"`
}

type DBConfig struct {
	User string `yaml:"db_user"`
	Pass string `yaml:"db_pass"`
	Host string `yaml:"db_host"`
	Port string `yaml:"db_port"`
	Name string `yaml:"db_name"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config path: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string = "config/config.yml"

	//flag.StringVar(&res, "config", "", "path to config file")
	//flag.Parse()
	//
	//if res == "" {
	//	res = os.Getenv("CONFIG_PATH")
	//}

	return res
}
