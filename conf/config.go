package conf

import "github.com/spf13/viper"

// AdminToken auth token for manage db
var AdminToken string

// Port app listen port
var Port string

// SqlitePath sqlite db path
var SqlitePath string

// OpenAIAPIBaseUrl OpenAI API base url
var OpenAIAPIBaseUrl string

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./data")

	viper.SetDefault("admin_token", "123456789")
	viper.SetDefault("port", "8080")
	viper.SetDefault("sqlite_path", "./data/oapi-proxy.db")
	viper.SetDefault("openai_api_base_url", "https://api.openai.com/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config file not found in ./data")
		} else {
			panic(err)
		}
	}

	AdminToken = viper.GetString("admin_token")
	iPort := viper.GetInt("port")
	SqlitePath = viper.GetString("sqlite_path")
	OpenAIAPIBaseUrl = viper.GetString("openai_api_base_url")

	if iPort > 0 && iPort < 65535 {
		Port = ":" + viper.GetString("port")
	}
}