package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	RefreshToken string
	QingTingId   string
}

func NewConfig() *Config {
	viper.SetConfigName(".qt-get")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	err := viper.ReadInConfig() // Find and read the config file
	config := &Config{}
	if err != nil { // Handle errors reading the config file
		config.RefreshToken = ""
		config.QingTingId = ""
		return config
	}
	config.RefreshToken = viper.GetString("refresh_token")
	config.QingTingId = viper.GetString("qingting_id")
	return config
}

func SaveConfig(config *Config) {
	viper.SetConfigName(".qt-get")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.Set("refresh_token", config.RefreshToken)
	viper.Set("qingting_id", config.QingTingId)
	err := viper.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error save config file: %w", err))
	}
}
