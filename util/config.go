package util

import "github.com/spf13/viper"

// config struct to hold the configurations params
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	StopFilePath  string `mapstructure:"STOP_FILE_PATH"`
}

//loads the config from the application env
func LoadConfig(path string) (config Config, err error) {

	// 1. Set the config name, type and add path
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	// 2. read the config
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// 3. unmarshal the configuration
	err = viper.Unmarshal(&config)
	return
}
