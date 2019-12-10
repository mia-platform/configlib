/*
 * Copyright 2019 Mia srl
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package configlib

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"gotest.tools/assert"
)

func TestGetConfigFromFile(t *testing.T) {
	type SubConfiguration struct {
		Kbool         bool                   `koanf:"kbool"`
		Kstring       string                 `koanf:"kstring"`
		Kint          int64                  `koanf:"kint"`
		Kfloat        float64                `koanf:"kfloat"`
		ArrayOfString []string               `koanf:"array-of-string"`
		ArrayOfNumber []int64                `koanf:"array-of-number"`
		Other         map[string]interface{} `koanf:"something"`
	}

	type Configuration map[string]SubConfiguration

	var config Configuration
	configName := "test-config.test"
	configPath := "."

	t.Run("read correctly configuration and set to config structure", func(t *testing.T) {
		err := GetConfigFromFile(configName, configPath, nil, &config)
		assert.Equal(t, err, nil, "Error is not nil.")
		assert.DeepEqual(t, config, Configuration{
			"lower-case-key": SubConfiguration{
				Kbool:         true,
				Kstring:       "my-string",
				Kint:          12,
				Kfloat:        float64(13.50),
				ArrayOfString: []string{"my", "values"},
				ArrayOfNumber: []int64{1, 2, 3},
				Other: map[string]interface{}{
					"a": "b",
					"c": "d",
				},
			},
			"UPPER-CASE-KEY": SubConfiguration{
				Kbool:         false,
				Kstring:       "my-other-string",
				Kint:          4,
				Kfloat:        float64(44.50),
				ArrayOfString: []string{"my", "other", "value"},
				ArrayOfNumber: []int64{6, 7, 8},
				Other: map[string]interface{}{
					"x": "y",
					"z": "v",
				},
			},
			"CamelCaseKey": SubConfiguration{
				Kint: 92,
			},
		})
	})

	t.Run("throw if config file not found at selected path", func(t *testing.T) {
		wrongConfigPath := "./wrong-path"
		err := GetConfigFromFile(configName, wrongConfigPath, nil, &config)
		assert.Assert(t, err != nil, "Error is nil.")
		_, ok := err.(viper.ConfigFileNotFoundError)
		assert.Assert(t, !ok, "Expected error is file not found, but not returned.")
	})

	t.Run("throw if config file not found with selected name", func(t *testing.T) {
		wrongConfigName := "a wrong file name"
		err := GetConfigFromFile(wrongConfigName, configPath, nil, &config)
		assert.Assert(t, err != nil, "Error is nil.")
		_, ok := err.(viper.ConfigFileNotFoundError)
		assert.Assert(t, !ok, "Expected error is file not found, but not returned.")
	})

	t.Run("throw if output struct does not contain config variable", func(t *testing.T) {
		type MyConfigStructure struct{}
		var wrongConfig MyConfigStructure
		err := GetConfigFromFile(configName, configPath, nil, &wrongConfig)
		assert.Assert(t, err != nil, "Error is nil.")
		assert.Equal(t, wrongConfig, MyConfigStructure{}, "Config is not empty.")
	})
}

func TestGetConfigFromFileWithJSONSchemaValidation(t *testing.T) {
	type SubConfiguration struct {
		Kbool         bool                   `koanf:"kbool"`
		Kstring       string                 `koanf:"kstring"`
		Kint          int64                  `koanf:"kint"`
		Kfloat        float64                `koanf:"kfloat"`
		ArrayOfString []string               `koanf:"array-of-string"`
		ArrayOfNumber []float64              `koanf:"array-of-number"`
		Other         map[string]interface{} `koanf:"something"`
	}

	type Configuration map[string]SubConfiguration

	var config Configuration
	configName := "test-config.test"
	configPath := "."

	t.Run("read correctly configuration, validate with json schema and set to config structure", func(t *testing.T) {
		jsonSchema := readFile(t, "./config.schema.test.json")
		err := GetConfigFromFile(configName, configPath, jsonSchema, &config)
		assert.Equal(t, err, nil, "Error is not nil.")
		assert.DeepEqual(t, config, Configuration{
			"lower-case-key": SubConfiguration{
				Kbool:         true,
				Kstring:       "my-string",
				Kint:          12,
				Kfloat:        float64(13.50),
				ArrayOfString: []string{"my", "values"},
				ArrayOfNumber: []float64{1, 2, 3},
				Other: map[string]interface{}{
					"a": "b",
					"c": "d",
				},
			},
			"UPPER-CASE-KEY": SubConfiguration{
				Kbool:         false,
				Kstring:       "my-other-string",
				Kint:          4,
				Kfloat:        float64(44.50),
				ArrayOfString: []string{"my", "other", "value"},
				ArrayOfNumber: []float64{6, 7, 8},
				Other: map[string]interface{}{
					"x": "y",
					"z": "v",
				},
			},
			"CamelCaseKey": SubConfiguration{
				Kint: 92,
			},
		})
	})

	t.Run("throws if config file validation fails", func(t *testing.T) {
		jsonSchema := []byte(`{
			"type": "string"
		}`)
		err := GetConfigFromFile(configName, configPath, jsonSchema, &config)
		assert.Assert(t, err != nil, "Error is not nil.")
		assert.Assert(t, strings.Contains(err.Error(), "configuration not valid:"), "Config file error.")
	})
}

func TestValidateJSONConfig(t *testing.T) {
	jsonSchema := []byte(`{
		"type": "object",
		"properties": {
			"object": {
				"type": "object",
				"properties": {
					"string": {
						"type": "string",
						"pattern": "^[\\w]+$"
					},
					"number": {"type": "number"},
					"boolean": {"type": "boolean"},
					"array": {
						"type": "array",
						"items": {"type": "number"}
					}
				}
			}
		}
	}`)

	t.Run("correctly validate json with json schema", func(t *testing.T) {
		document := `{
			"object": {
				"string": "my_string",
				"number": 2,
				"boolean": true,
				"array": [0, 5, 10, 15]
			}
		}`

		err := validateJSONConfig(jsonSchema, []byte(document))
		assert.Equal(t, err, nil, "JSON document validation returns error.")
	})

	t.Run("throws if json schema is not correct", func(t *testing.T) {
		document := []byte(``)
		wrongSchema := []byte(`{"type": "not-correct-type"}`)
		err := validateJSONConfig(wrongSchema, document)
		assert.Assert(t, err != nil, "JSON document should returns error.")
		assert.Assert(t, strings.Contains(err.Error(), "error validating:"), "", "JSON document should returns the correct error.")
	})

	t.Run("throws if document is not valid", func(t *testing.T) {
		document := `{
			"object": {
				"string": "my_string",
				"number": "2",
				"boolean": true,
				"array": [0, 5, 10, 15]
			}
		}`
		err := validateJSONConfig(jsonSchema, []byte(document))
		assert.Assert(t, err != nil, "JSON document should returns error.")
		assert.Assert(t, strings.Contains(err.Error(), "json schema validation errors:"), "", "JSON document should returns the correct error.")
	})
}

func closeJSONFile(t *testing.T, jsonFile *os.File) {
	t.Helper()
	err := jsonFile.Close()
	assert.Equal(t, err, nil, "Error setting env variables")
}

func readFile(t *testing.T, filePath string) []byte {
	t.Helper()
	jsonFile, err := os.Open(filePath)
	defer closeJSONFile(t, jsonFile)
	assert.Equal(t, err, nil, "Failed to open json schema file.")
	byteValue, err := ioutil.ReadAll(jsonFile)
	assert.Equal(t, err, nil, "Failed to read json schema file.")
	return byteValue
}
