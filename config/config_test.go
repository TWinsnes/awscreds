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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretBackendAllowedValues(t *testing.T) {
	conf := config{}

	expected := []string{"aws"}
	actual := conf.SecretBackendAllowedValues()

	assert.Equal(t, expected, actual)
}

func TestSetSecretBackend(t *testing.T) {
	conf := config{}

	allowed := conf.SecretBackendAllowedValues()

	for _, value := range allowed {
		err := conf.SetSecretBackend(value)

		assert.NoError(t, err, "Expected '"+value+"' to be allowed")
	}

	err := conf.SetSecretBackend("notallowedvalue")

	assert.Error(t, err, "Expected 'notallowedvalue' to produce error")

}

func TestSetSecretBackendCapInsensitive(t *testing.T) {
	conf := config{}

	err := conf.SetSecretBackend("aWs")

	assert.NoError(t, err, "Expected 'aWs' to be allowed")
}

func TestValidateSecretBackend(t *testing.T) {
	conf := config{}
	conf.loadDefaults()

	allowed := conf.SecretBackendAllowedValues()

	for _, value := range allowed {
		conf.wrapper.SecretBackend = value

		err := conf.Validate()

		assert.NoError(t, err, "Expected '"+value+"' to be allowed")
	}

	conf.wrapper.SecretBackend = "notvalid"

	err := conf.Validate()

	assert.Error(t, err, "Expected 'notvalid' to fail validation")

}
