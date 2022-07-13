package ConfigFuncs

import (
	"text/template"

	"github.com/spf13/viper"
)

var secured = "n" // "n" (no) will merge insecure and secure functions unless stated by ldflags

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
	txtFuncSMap := template.FuncMap{
		"getConfigString":     config.GetConfigString,
		"askGetConfigString":  config.AskGetConfigString,
		"mustGetConfigString": config.MustGetConfigString,
	}

	txtFuncIMap := template.FuncMap{}

	if secured == "n" {
		// iterate through insecure map and append to secure map
		for k, v := range txtFuncIMap {
			txtFuncSMap[k] = v
		}
	}

	return txtFuncSMap
}
