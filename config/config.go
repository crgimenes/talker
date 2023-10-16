package config

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
	Salt          string `ini:"salt" cfg:"salt" cfgRequired:"true" cfgHelper:"Salt used to hash passwords."`
	Listen        string `ini:"listen" cfg:"listen" cfgDefault:"0.0.0.0:2200"`
	DatabaseURL   string `ini:"database_url" cfg:"database_url" cfgRequired:"true"`
	FilesPath     string `ini:"files_path" cfg:"files_path" cfgRequired:"true"`
	SecureCookie  bool   `ini:"secure_cookie" cfg:"secure_cookie" cfgDefault:"false" cfgHelper:"If true, the cookie will only be sent over HTTPS connections, never over plain HTTP."`
	MaxAgeSession int    `ini:"max_age_session" cfg:"max_age_session" cfgDefault:"86400" cfgHelper:"Max age of session in seconds."`
	MaxUploadSize int64  `ini:"max_upload_size" cfg:"max_upload_size" cfgDefault:"10485760" cfgHelper:"Max upload size in bytes."`
}

var CFG *Config

func Load() error {
	CFG = &Config{}
	config.PrefixEnv = "TALKER"
	config.File = "talker.ini"
	return config.Parse(CFG)
}
