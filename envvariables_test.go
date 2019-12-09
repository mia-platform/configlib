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

package helpers

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"gotest.tools/assert"
)

func testSetEnv(t *testing.T, key, value string) {
	t.Helper()
	err := os.Setenv(key, value)
	assert.Equal(t, err, nil, "Error setting env variables")
}

func testUnsetEnv(t *testing.T, key string) {
	t.Helper()
	err := os.Unsetenv(key)
	assert.Equal(t, err, nil, "Error setting env variables")
}

func TestGetEnvVariables(t *testing.T) {
	t.Run("returns env variables", func(t *testing.T) {
		type MyEnvDef struct {
			MyDefault string
			Custom    int
			Required  bool
		}

		var env MyEnvDef
		envConfig := []EnvConfig{
			{
				Key:          "MY_DEFAULT_ENV",
				Variable:     "MyDefault",
				DefaultValue: "default-value",
			},
			{
				Key:          "MY_CUSTOM_ENV",
				Variable:     "Custom",
				DefaultValue: "5",
			},
			{
				Key:      "MY_REQUIRED_ENV",
				Variable: "Required",
				Required: true,
			},
		}

		testSetEnv(t, "MY_CUSTOM_ENV", "10")
		defer testUnsetEnv(t, "MY_CUSTOM_ENV")
		testSetEnv(t, "MY_REQUIRED_ENV", "true")
		defer testUnsetEnv(t, "MY_REQUIRED_ENV")

		err := GetEnvVariables(envConfig, &env)

		assert.Equal(t, err, nil, "Error getting values.")
		assert.Equal(t, env.MyDefault, envConfig[0].DefaultValue, "Not returns default value.")
		assert.Equal(t, env.Custom, 10, "Not returns custom value.")
		assert.Equal(t, env.Required, true, "Not returns required value.")
	})

	t.Run("returns env variables with default", func(t *testing.T) {
		type MyEnvDef struct {
			String  string
			Integer int
			Boolean bool
		}

		var env MyEnvDef
		envConfig := []EnvConfig{
			{
				Key:          "MY_DEFAULT_ENV",
				Variable:     "String",
				DefaultValue: "default-value",
			},
			{
				Key:          "MY_INT_ENV",
				Variable:     "Integer",
				DefaultValue: "5",
			},
			{
				Key:          "MY_BOOLEAN_ENV",
				Variable:     "Boolean",
				DefaultValue: "true",
			},
		}

		err := GetEnvVariables(envConfig, &env)

		assert.Equal(t, err, nil, "Error getting values.")
		assert.Equal(t, env.String, "default-value", "Not returns default value.")
		assert.Equal(t, env.Integer, 5, "Not returns correct default value.")
		assert.Equal(t, env.Boolean, true, "Not returns correct default value.")
	})

	t.Run("returns error if required value not set", func(t *testing.T) {
		type MyEnvDef struct {
			MyEnv string
		}

		var env MyEnvDef
		const customValue string = "my-custom-value"
		envConfig := []EnvConfig{
			{
				Key:      "MY_ENV",
				Variable: "MyEnv",
				Required: true,
			},
		}

		err := GetEnvVariables(envConfig, &env)

		assert.Assert(t, err != nil, "Get env variables not errored.")
		assert.Equal(t, err.Error(), "required env variable MY_ENV not set", "Get env variables not errored.")
		assert.Equal(t, env, MyEnvDef{}, "Not returns value.")
	})

	t.Run("throw if output struct does not contain config variable", func(t *testing.T) {
		type MyEnvDef struct{}

		var env MyEnvDef
		const customValue string = "my-custom-value"
		envConfig := []EnvConfig{
			{
				Key:          "MY_ENV",
				Variable:     "MyEnv",
				DefaultValue: "default-value",
			},
		}

		err := GetEnvVariables(envConfig, &env)

		assert.Assert(t, err != nil, "Get env variables not errored.")
		assert.Equal(t, env, MyEnvDef{}, "Not returns value.")
	})

	t.Run("not throws if env config contains empty variable and structure is empty", func(t *testing.T) {
		type MyEnvDef struct{}

		var env MyEnvDef
		const customValue string = "my-custom-value"
		envConfig := []EnvConfig{
			{
				Key:      "MY_ENV",
				Variable: "MyEnv",
			},
		}

		err := GetEnvVariables(envConfig, &env)

		assert.Equal(t, err, nil, "Get env variables not errored.")
		assert.Equal(t, env, MyEnvDef{}, "Not returns value.")
	})

	t.Run("not throws if empty struct output is not present into env config", func(t *testing.T) {
		type MyEnvDef struct {
			MyEnv string
		}

		var env MyEnvDef
		const customValue string = "my-custom-value"
		envConfig := []EnvConfig{}

		err := GetEnvVariables(envConfig, &env)

		assert.Equal(t, err, nil, "Get env variables not errored.")
		assert.Equal(t, env, MyEnvDef{}, "Not returns value.")
	})
}

func TestSetViperVariables(t *testing.T) {
	t.Run("if env is required but not given, throw error", func(t *testing.T) {
		v := viper.New()
		err := setViperVariable(v, EnvConfig{
			Key:      "MY_REQUIRED_TEST_KEY",
			Variable: "variableKey",
			Required: true,
		})
		assert.Assert(t, err != nil, "Set viper variable not throw if required variable does not exist.")
		assert.Equal(t, err.Error(), "required env variable MY_REQUIRED_TEST_KEY not set", "Set variable error message is not correct")
	})

	t.Run("if env has default but has not set value, set to viper the default value", func(t *testing.T) {
		v := viper.New()
		const key = "ENV_KEY"
		const defaultValue = "value"
		const variableKey = "variableKey"
		testUnsetEnv(t, key)
		err := setViperVariable(v, EnvConfig{
			Key:          key,
			Variable:     variableKey,
			DefaultValue: defaultValue,
		})
		assert.Equal(t, err, nil, "set env variable does not return error")
		assert.Equal(t, v.GetString(variableKey), defaultValue, "Get a not correct variable.")
	})

	t.Run("if env has value, set value", func(t *testing.T) {
		v := viper.New()
		const variableKey = "variableKey"
		const key = "MY_REQUIRED_TEST_KEY"
		const customValue = "my-custom-value"
		testSetEnv(t, key, customValue)
		defer testUnsetEnv(t, key)
		err := setViperVariable(v, EnvConfig{
			Key:          key,
			Variable:     variableKey,
			DefaultValue: "value",
		})
		assert.Equal(t, err, nil, "set env variable does not return error")
		assert.Equal(t, v.GetString(variableKey), customValue, "Get a not correct variable.")
	})
}
