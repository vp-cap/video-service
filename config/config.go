package config

import (
	"log"

	dataConfig "github.com/vp-cap/data-lib/config"

	viper "github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Services ServiceConfigurations
	Database dataConfig.DatabaseConfiguration
	Storage  dataConfig.StorageConfiguration
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port string
}

// ServiceConfigurations exported
type ServiceConfigurations struct {
	TaskAllocator string
}

// GetConfigs Get Configurations from config.yaml and set in Configurations struct
func GetConfigs() (Configurations, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // type
	viper.AddConfigPath("config") // optionally look for config in the working directory
	viper.AutomaticEnv()          // enable viper to read env

	// store in configuration struct
	var configs Configurations

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return configs, err
	}
	if err := viper.Unmarshal(&configs); err != nil {
		log.Println(err)
		return configs, err
	}
	return configs, nil
}