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
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// EnvConfig to setup to access to env variables
type EnvConfig struct {
	Key          string
	Variable     string
	DefaultValue string
	Required     bool
}

// GetEnvVariables env variables
func GetEnvVariables(envVariablesConfig []EnvConfig, output interface{}) error {
	v := viper.New()
	v.AutomaticEnv()
	for _, config := range envVariablesConfig {
		if err := setViperVariable(v, config); err != nil {
			return err
		}
	}

	if err := v.UnmarshalExact(&output); err != nil {
		return fmt.Errorf("unable to decode into struct: %s", err.Error())
	}
	return nil
}

func setViperVariable(v *viper.Viper, env EnvConfig) error {
	if err := v.BindEnv(env.Variable, env.Key); err != nil {
		return fmt.Errorf("Error binding env variable %s", env.Key)
	}
	if env.DefaultValue != "" {
		v.SetDefault(env.Variable, env.DefaultValue)
	}
	if env.Required {
		_, ok := os.LookupEnv(env.Key)
		if !ok {
			return fmt.Errorf("required env variable %s not set", env.Key)
		}
	}
	return nil
}
