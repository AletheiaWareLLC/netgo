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

package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logs := os.Args[1:]
	if len(logs) == 0 {
		logs = append(logs, "logs")
	}

	os.Remove("log.db")

	count, err := Parse("log.db", []string{
		"log.go:",
		"net.go:39:",
		"net.go:44:",
		"net.go:47:",
	}, logs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Parsed", count, "Records")
}
