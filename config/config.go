package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// Config contains configuration values for util
type Config struct {
	SecretBackend string `toml:"SecretBackend"`
}

// LoadConfig loads the configuration from file,
// and falls back to default calues if file
// could not be be loaded
func LoadConfig(path string) Config {
	var conf = Config{}
	err := conf.loadConfigFile(path)

	if err != nil {
		conf.loadDefaults()
	}

	return conf
}

func (c *Config) loadDefaults() {
	c.SecretBackend = "aws"
}

// LoadConfigFile loads configuration from a toml file
func (c *Config) loadConfigFile(path string) error {

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	tomlcontent := string(bytes)

	toml.Decode(tomlcontent, c)

	return nil
}

// SaveConfig saves config to toml file
func (c *Config) SaveConfig(path string) error {

	writer, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return nil
	}

	defer writer.Close()

	encoder := toml.NewEncoder(writer)

	encoder.Encode(c)

	return nil
}
