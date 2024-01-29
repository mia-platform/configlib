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
	"encoding/json"
	"fmt"

	kJson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
)

func validateJSONConfig(schema, jsonConfig []byte) error {
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(jsonConfig)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("error validating: %s", err.Error())
	}
	if !result.Valid() {
		return fmt.Errorf("json schema validation errors: %s", result.Errors())
	}
	return nil
}

// GetConfigFromFile func read configuration from file and save in output interface.
func GetConfigFromFile(configName, configPath string, jsonSchema []byte, output interface{}) error {
	var k = koanf.New(".")

	if err := k.Load(file.Provider(fmt.Sprintf("%s/%s.json", configPath, configName)), kJson.Parser()); err != nil {
		return fmt.Errorf("error loading config file: %s", err.Error())
	}

	if jsonSchema != nil {
		jsonDocument, err := json.Marshal(k.Raw())
		if err != nil {
			return fmt.Errorf("config document stringify failed: %s", err.Error())
		}
		err = validateJSONConfig(jsonSchema, jsonDocument)
		if err != nil {
			return fmt.Errorf("configuration not valid: %s", err.Error())
		}
	}

	if err := k.UnmarshalWithConf("", &output, koanf.UnmarshalConf{
		DecoderConfig: &mapstructure.DecoderConfig{
			Metadata:         nil,
			Result:           &output,
			WeaklyTypedInput: true,
			ErrorUnused:      true,
		},
	}); err != nil {
		return fmt.Errorf("error unmarshalling file: %s", err.Error())
	}
	return nil
}
