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
	"fmt"
	"io/ioutil"
	"os"
)

// ReadFile is a utility to read a file from the file system
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file error: %s", err.Error())
	}
	defer func(fileToClose *os.File) {
		err := fileToClose.Close()
		if err != nil {
			panic(fmt.Errorf("error closing file %s: %s", filePath, err.Error()))
		}
	}(file)
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %s", filePath, err.Error())
	}
	return byteValue, nil
}
