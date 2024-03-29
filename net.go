/*
 * Copyright 2019 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package netgo

import (
	"net/url"
	"os"
	"strconv"
)

func BooleanFlag(name string) bool {
	flag, ok := os.LookupEnv(name)
	if !ok {
		return false
	}
	b, err := strconv.ParseBool(flag)
	if err != nil {
		return false
	}
	return b
}

func QueryParameter(query url.Values, parameter string) string {
	if results, ok := query[parameter]; ok && len(results) > 0 {
		return results[0]
	}
	return ""
}
