/*
 * Copyright 2022 Aletheia Ware LLC
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
    "strings"
    "strconv"
    "log"
)

func ParseInt(s string) int64 {
    s = strings.TrimSpace(s)
    if s != "" {
        if i, err := strconv.ParseInt(s, 10, 64); err != nil {
            log.Println(err)
        } else {
            return int64(i)
        }
    }
    return 0
}