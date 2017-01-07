package evedt

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

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

// Config
type Config struct {
	ApiCredentials ApiCredentials   `toml:"eveapi"`
	Replicator     ReplicatorConfig `toml:"replicator"`
	ServerConfig   ServerConfig     `toml:"server"`
}

func decodeConfig(tml string) (*Config, error) {

	conf := &Config{}
	if _, err := toml.Decode(tml, conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// LoadConfig loads and parses the specified config file
func LoadConfig(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	str := string(bytes)
	conf, err := decodeConfig(str)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
