package evedt

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// decodeConfig parses a string into it's struct representation
func decodeConfig(tml string, conf interface{}) error {

	if _, err := toml.Decode(tml, conf); err != nil {
		return err
	}
	return nil
}

// LoadConfig loads and parses the specified filename into an interface, most likely a user defined *struct.
// In case anything goes wrong it will return a non-nil error.
func LoadConfig(filename string, conf interface{}) error {

	var err error
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	str := string(bytes)
	err = decodeConfig(str, conf)
	if err != nil {
		return err
	}
	return nil
}
