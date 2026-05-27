package settings

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App          App          `yaml:"app"`
	Database     Database     `yaml:"database"`
	Apikeys      ApiKeys      `yaml:"api_keys"`
	AppleKeys    AppleKeys    `yaml:"apple_keys"`
	EmailService EmailService `yaml:"email_service"`
	PushService  PushService  `yaml:"push_service"`
}

type App struct {
	Port      int    `yaml:"port"`
	Debug     bool   `yaml:"debug"`
	LogLevel  string `yaml:"log_level"`
	SecretKey string `yaml:"secret_key"`
	Env       string `yaml:"env"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DataBase string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type ApiKeys struct {
	OpenAI string `yaml:"openai"`
}

type AppleKeys struct {
	AppleIapKeyId          string `yaml:"apple_iap_key_id"`
	AppleIapIssuerId       string `yaml:"apple_iap_issuer_id"`
	AppleBundleId          string `yaml:"apple_bundle_id"`
	AppleIapPrivateKeyPath string `yaml:"apple_iap_private_key_path"`
	AppleIapBaseUrl        string `yaml:"apple_iap_base_url"`
	AppleIapSandboxBaseUrl string `yaml:"apple_iap_sandbox_base_url"`
}

type EmailService struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type PushService struct {
	APNsKeyID          string `yaml:"apns_key_id"`
	APNsTeamID         string `yaml:"apns_team_id"`
	APNsTopic          string `yaml:"apns_topic"`
	APNsPrivateKeyPath string `yaml:"apns_private_key_path"`
	APNsBaseURL        string `yaml:"apns_base_url"`
	APNsSandboxBaseURL string `yaml:"apns_sandbox_base_url"`
	UseSandbox         bool   `yaml:"use_sandbox"`
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
