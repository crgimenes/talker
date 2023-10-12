package config

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
}

var CFG *Config

func Load() error {
	CFG = &Config{}
	config.PrefixEnv = "TALKER"
	config.File = "talker.ini"
	return config.Parse(CFG)
}
