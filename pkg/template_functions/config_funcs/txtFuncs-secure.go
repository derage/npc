package ConfigFuncs

import "fmt"

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
		return "", fmt.Errorf("value %s not set", key)
	}
	return config.GetConfigString(key), nil
}
