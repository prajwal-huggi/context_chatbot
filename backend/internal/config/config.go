package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct{
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct{
	Env string `yaml:"env" env:"ENV" env-required:"true" env-default:"Prod"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config{
	var configPath string
	
	configPath= os.Getenv("CONFIG_PATH")

	if configPath==""{
		flags:= flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath= *flags

		if configPath== ""{
			log.Fatal("Config path is not set")
		}
	}

	if _, err:= os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var config Config

	err:= cleanenv.ReadConfig(configPath, &config)

	if err!= nil{
		log.Fatal("Cannot read the config file")
	}

	return &config
}
