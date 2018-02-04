// Copyright © 2018 Thomas Winsnes <twinsnes@live.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Handler handles loading, storing, validating the configuration file
type Handler interface {
	// config value getters and setters
	SecretBackend() string
	SecretBackendAllowedValues() []string
	SetSecretBackend(value string) error

	// handler functions
	SaveConfig(path string) error
	Validate() error
}

// internalConfig wraps the config values as the toml library was
// having issue with getters and setters on the struct
type internalConfig struct {
	SecretBackend string `toml:"SecretBackend"`
}

// Config contains configuration values for util
type config struct {
	wrapper internalConfig
}

// LoadConfig loads the configuration from file,
// and falls back to default calues if file
// could not be be loaded
func LoadConfig(path string) Handler {
	var conf = config{}
	err := conf.loadConfigFile(path)

	if err != nil {
		conf.loadDefaults()
	}

	return &conf
}

func (c *config) SecretBackend() string {
	return c.wrapper.SecretBackend
}

func (c *config) SetSecretBackend(value string) error {

	transformedValue := strings.ToLower(value)

	valid := false
	for _, allowedValue := range c.SecretBackendAllowedValues() {
		if allowedValue == transformedValue {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("SetSecretBackend does not allow the value '%s'. \nAllowed values: %s", value, c.SecretBackendAllowedValues())
	}

	c.wrapper.SecretBackend = transformedValue

	return nil
}

func (c *config) SecretBackendAllowedValues() []string {
	return []string{"aws"}
}

func (c *config) loadDefaults() {
	c.wrapper.SecretBackend = "aws"
}

// LoadConfigFile loads configuration from a toml file
func (c *config) loadConfigFile(path string) error {

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	tomlcontent := string(bytes)

	fmt.Println("TOML content: " + tomlcontent)

	wrapper := internalConfig{}

	_, err = toml.Decode(tomlcontent, wrapper)

	if err != nil {
		return err
	}

	c.wrapper = wrapper

	return nil
}

// SaveConfig saves config to toml file
func (c *config) SaveConfig(path string) error {

	writer, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return nil
	}

	defer writer.Close()

	encoder := toml.NewEncoder(writer)

	return encoder.Encode(c.wrapper)
}

// Validate validates the loaded configuration
func (c *config) Validate() error {

	return c.SetSecretBackend(c.wrapper.SecretBackend)
}
