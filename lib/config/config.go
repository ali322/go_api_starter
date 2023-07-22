package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var App = new(AppConf)

type AppConf struct {
	Port       int64  `mapstructure:"port"`
	LogDir     string `mapstructure:"logDir"`
	RecordDir  string `mapstructure:"recordDir"`
	RecordPath string `mapstructure:"recordPath"`
}

func Read() {
	workDir, _ := os.Getwd()
	viper.SetConfigFile(filepath.Join(workDir, "app.toml"))
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Sub("app").Unmarshal(App); err != nil {
		log.Fatal(err)
	}
}
