package ConfigFuncs

import (
	"fmt"
	"text/template"

	"github.com/spf13/viper"
)

type Config struct {
	Viper *viper.Viper
}

func InitializeConfigFuncs(viperConfig *viper.Viper) template.FuncMap {
	configFunctions := Config{
		Viper: viperConfig,
	}

	return configFunctions.TxtFuncMap()
}

func (config *Config) TxtFuncMap() template.FuncMap {
	configFuncMap := template.FuncMap{
		"getConfigString":     config.GetConfigString,
		"askGetConfigString":  config.AskGetConfigString,
		"mustGetConfigString": config.MustGetConfigString,
	}

	return configFuncMap
}

func (config *Config) GetConfigString(key string) string {
	return config.Viper.GetString(key)
}

func (config *Config) AskGetConfigString(key string) string {
	if !config.Viper.IsSet(key) {
		var retrivedValue string
		fmt.Println("Enter value for " + key + ": ")
		fmt.Scanln(&retrivedValue)
		config.Viper.Set(key, retrivedValue)
	}
	return config.GetConfigString(key)
}

func (config *Config) MustGetConfigString(key string) (string, error) {
	if !config.Viper.IsSet(key) {
		return "", fmt.Errorf("Value %s not set", key)
	}
	return config.GetConfigString(key), nil
}
