package evedt

import tconf "github.com/trichner/evedt/tconf"

type ApiCredentials struct {
	ApiKey string `toml:"api-key"`
	VCode  string `toml:"verification-code"`
}

type ReplicatorConfig struct {
	Interval int `toml:"interval"`
}

type ServerConfig struct {
	Port   string `toml:"port"`
	Prefix string `toml:"prefix"`
}

// Config struct holds the content of the configuration file
type Config struct {
	ApiCredentials ApiCredentials   `toml:"eveapi"`
	Replicator     ReplicatorConfig `toml:"replicator"`
	ServerConfig   ServerConfig     `toml:"server"`
}

func loadConfig(filename string) (*Config, error) {

	conf := &Config{}
	err := tconf.LoadConfig(filename, conf)
	return conf, err
}
