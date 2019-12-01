package config

import (
	"ImgCrop/internal/structs"
	"github.com/spf13/viper"
)

func GetConfig() (structs.Config, error) {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath("./config") // path to look for the config file in
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return structs.Config{}, err
	}

	var config structs.Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return structs.Config{}, err
	}

	return config, nil
}
