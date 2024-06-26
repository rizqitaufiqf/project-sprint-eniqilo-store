package config

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfigs *envConfigs

func init() {
	EnvConfigs = loadENV()
}

type envConfigs struct {
	DbName     string `mapstructure:"DB_NAME"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbUser     string `mapstructure:"DB_USERNAME"`
	DbParams   string `mapstructure:"DB_Params"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbPort     string `mapstructure:"DB_PORT"`
	BcryptSalt string `mapstructure:"BCRYPT_SALT"`
	JwtSecret  string `mapstructure:"JWT_SECRET"`
}

func loadENV() (config *envConfigs) {
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return config
}
