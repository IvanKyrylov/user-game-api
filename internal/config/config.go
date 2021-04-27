package config

import (
	"sync"

	"github.com/IvanKyrylov/user-game-api/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	MongoDB struct {
		Host                string `yaml:"host" env-required:"true"`
		Port                string `yaml:"port" env-required:"true"`
		Username            string `yaml:"username" env-required:"true"`
		Password            string `yaml:"password" env-required:"true"`
		AuthDB              string `yaml:"auth_db" env-required:"true"`
		Database            string `yaml:"database" env-required:"true"`
		CollectionUsers     string `yaml:"collection_users" env-required:"true"`
		CollectionUserGames string `yaml:"collection_user_games" env-required:"true"`
	} `yaml:"mongodb" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logging.CommonLog.Println("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logging.CommonLog.Println(help)
			logging.ErrorLog.Fatal(err)
		}
	})
	return instance
}
