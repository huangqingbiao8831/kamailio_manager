package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config reflects the structure of managerConf.yaml
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Kamailio struct {
		RPCUrl  string `mapstructure:"rpc_url"`
		Timeout int    `mapstructure:"timeout"`
	} `mapstructure:"kamailio"`
	Supervisor struct {
		Url         string `mapstructure:"url"`
		ProcessName string `mapstructure:"process_name"`
	} `mapstructure:"supervisor"`
}

// LoadConfig initializes viper and unmarshals the configuration
func LoadConfig() *Config {
	setDefaults()

	viper.SetConfigName("managerConf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/usr/local/kamailio/etc/kamailio") // Config file path

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults")
		} else {
			log.Fatalf("Read config file failed: %v", err)
		}
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Unmarshal config failed: %v", err)
	}
	return &conf
}

func setDefaults() {
	viper.SetDefault("server.port", "8083")
	viper.SetDefault("kamailio.rpc_url", "http://127.0.0.1:8091/rpc")
	viper.SetDefault("kamailio.timeout", 5)
	viper.SetDefault("supervisor.url", "/var/run/supervisor.sock")
	viper.SetDefault("supervisor.process_name", "kamailio")
}
