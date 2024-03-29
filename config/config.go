package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

func LoadTestConfig(path string) Config {
	vip := viper.New()
	vip.SetConfigName("config_test")
	vip.SetConfigType("yml")
	vip.AddConfigPath(path)
	return loadConfigWithViper(vip)
}

func Load() Config {
	vip := viper.New()

	vip.SetConfigName("config")
	vip.SetConfigType("yml")
	vip.AddConfigPath(".")

	return loadConfigWithViper(vip)
}

func loadConfigWithViper(vip *viper.Viper) Config {
	vip.SetEnvPrefix("docker")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vip.AutomaticEnv()

	err := vip.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// workaround https://github.com/spf13/viper/issues/188#issuecomment-399518663
	// to allow read from environment variables when Unmarshal
	for _, key := range vip.AllKeys() {
		val := vip.Get(key)
		vip.Set(key, val)
	}

	fmt.Println("Config file used:", vip.ConfigFileUsed())

	cfg := Config{}
	err = vip.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
