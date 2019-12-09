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
	"strings"
	"testing"

	"gotest.tools/assert"
)

func stripStrings(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(string(s))), " ")
}

func TestReadFile(t *testing.T) {
	t.Run("Correctly read file", func(t *testing.T) {
		file, err := ReadFile("./test-config.test.json")
		assert.Equal(t, err, nil, "Error reading existent file.")
		assert.Equal(t, stripStrings(string(file)), stripStrings(`{
			"lower-case-key": {
				"kbool": true,
				"kstring": "my-string",
				"kint": 12,
				"kfloat": 13.50,
				"array-of-string": [
					"my",
					"values"
				],
				"array-of-number": [
					1,
					2,
					3
				],
				"something": {
					"a": "b",
					"c": "d"
				}
			},
			"UPPER-CASE-KEY": {
				"kbool": false,
				"kstring": "my-other-string",
				"kint": 4,
				"kfloat": 44.50,
				"array-of-string": [
					"my",
					"other",
					"value"
				],
				"array-of-number": [
					6,
					7,
					8
				],
				"something": {
					"x": "y",
					"z": "v"
				}
			},
			"CamelCaseKey": {
				"kint": 92
			}
		}`))
	})

	t.Run("Throws if file not exists", func(t *testing.T) {
		_, err := ReadFile("./not-real-file")
		assert.Assert(t, err != nil, "Error reading existent file.")
		assert.Assert(t, strings.Contains(err.Error(), "open file error:"))
	})
}
