package config

import (
	"om/pkg/util"
	"os"
)

type Config struct {
	Listen    string `json:"listen"`
	SqlDriver string `json:"sql_driver"`
	DSN       string `json:"dsn"`
	JwtKey    string `json:"jwt_key"`
	LogLevel  string `json:"log_level"`
}

var Cfg Config

func LoadOrCreate() error {
	cfgPath := util.DataDir + "config.json"
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		Cfg = Config{
			Listen:    ":34567",
			SqlDriver: "sqlite",
			DSN:       util.DataDir + "data.db?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-16000)",
			JwtKey:    util.RandStr(32),
			LogLevel:  "error",
		}

		err = Save(&Cfg)
		if err != nil {
			return err
		}

		return nil
	}

	err = util.Json.Unmarshal(data, &Cfg)
	if err != nil {
		return err
	}
	return nil
}

func Save(cfg *Config) error {
	data, err := util.Json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(util.DataDir+"config.json", data, 0600)
}
