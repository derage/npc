package ConfigFuncs

import (
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
		"getConfigString": config.GetConfigString,
	}

	return configFuncMap
}

func (config *Config) GetConfigString(key string) string {
	return config.Viper.GetString(key)
}
