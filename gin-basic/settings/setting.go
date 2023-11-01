package settings

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DataBase string `yaml:"database"`
		Charset  string `yaml:"charset"`
	} `yaml:"database"`

	App struct {
		Port      int    `yaml:"port"`
		Debug     bool   `yaml:"debug"`
		LogLevel  string `yaml:"log_level"`
		SecretKey string `yaml:"secret_key"`
		Env       string `yaml:"env"`
	} `yaml:"app"`
	ApiKeys struct {
		Google   string `yaml:"google"`
		Facebook string `yaml:"facebook"`
	} `yaml:"api_keys"`
}

// 配置开发环境
var env = getEnvironment()
var Conf Config

func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		// 如果未设置 ENVIRONMENT 环境变量，默认为开发环境
		return "dev"
	}
	return env
}

func LoadConfig() error {
	configFile := "config/dev.yaml" // 默认使用开发环境配置
	if env == "production" {
		configFile = "/config/production.yaml"
	} else if env == "dev" {
		configFile = "/config/dev.yaml"
	} else {
		configFile = "config/local.yaml"
	}
	// 读取配置文件
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		// 处理错误
		return err
	}

	// 解析配置
	var configInfo Config
	if err := yaml.Unmarshal(data, &configInfo); err != nil {
		// 处理错误
		return err
	}
	Conf = configInfo
	return nil
}
