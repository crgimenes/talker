package config

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
	Listen      string `ini:"listen" cfg:"listen" cfgDefault:"0.0.0.0:2200"`
	DatabaseURL string `ini:"database_url" cfg:"database_url" cfgRequired:"true"`
	FilesPath   string `ini:"files_path" cfg:"files_path" cfgRequired:"true"`
}

var CFG *Config

func Load() error {
	CFG = &Config{}
	config.PrefixEnv = "TALKER"
	config.File = "talker.ini"
	return config.Parse(CFG)
}
