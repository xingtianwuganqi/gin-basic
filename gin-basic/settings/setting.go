package settings

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	App App  `yaml:"app"`
	Log Log  `yaml:"log"`
}

type App struct {
	Port      int    `yaml:"port"`
	Debug     bool   `yaml:"debug"`
	LogLevel  string `yaml:"log_level"`
	SecretKey string `yaml:"secret_key"`
	Env       string `yaml:"env"`
}

type Log struct {
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

var Conf = &Config{}

func InitConf() {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "local" // 默认环境
	}

	configPath := "config/" + environment + ".yaml"

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic("Failed to read config file: " + configPath + ", error: " + err.Error())
	}

	err = yaml.Unmarshal(data, Conf)
	if err != nil {
		panic("Failed to parse config file, error: " + err.Error())
	}
}