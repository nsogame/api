package api

import "github.com/spf13/viper"

type Config struct {
	SecretKey string
	Debug     bool
	BindAddr  string

	DbProvider   string
	DbConnection string
}

func GetConfig() (config Config, err error) {
	v := viper.New()
	v.SetConfigName("api")

	v.SetDefault("SecretKey", "hello")
	v.SetDefault("Debug", "false")
	v.SetDefault("BindAddr", "127.0.0.1:8000")
	v.SetDefault("DbProvider", "sqlite3")
	v.SetDefault("DbConnection", "api.db")

	v.AddConfigPath(".")
	err = v.ReadInConfig()
	if err != nil {
		return
	}

	err = v.Unmarshal(&config)
	return
}
