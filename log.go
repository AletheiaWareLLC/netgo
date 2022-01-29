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
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SetupLogging() (*os.File, error) {
	store, ok := os.LookupEnv("LOG_DIRECTORY")
	if !ok {
		store = "logs"
	}
	if err := os.MkdirAll(store, os.ModePerm); err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile(filepath.Join(store, time.Now().UTC().Format(time.RFC3339)), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return logFile, nil
}

func LogRequest(r *http.Request) {
	log.Println(r.RemoteAddr, r.Proto, r.Method, r.Host, r.URL, r.Header)
}

func IsRequestLog(sources []string, line string) bool {
	if !strings.HasPrefix(line, "2") {
		return false
	}
	if _, err := time.Parse("2006/01/02 15:04:05", line[:19]); err != nil {
		return false
	}
	for _, s := range sources {
		if strings.HasPrefix(line[20:], s) {
			return true
		}
	}
	return false
}

func ParseRequestLog(line string) (int64, []string, map[string]string, error) {
	var (
		source,
		address,
		protocol,
		method,
		host,
		url string
		headers map[string]string
	)

	start := 20
	timestamp, err := time.Parse("2006/01/02 15:04:05 ", line[:start])
	if err != nil {
		return 0, nil, nil, err
	}

	end := start + strings.IndexRune(line[start:], ' ')
	source = line[start:end]

	start = end + 1
	end = start + strings.IndexRune(line[start:], ' ')
	address, _, err = net.SplitHostPort(line[start:end])
	if err != nil {
		return 0, nil, nil, err
	}

	start = end + 1
	end = start + strings.IndexRune(line[start:], ' ')
	protocol = line[start:end]

	start = end + 1
	end = start + strings.IndexRune(line[start:], ' ')
	method = line[start:end]

	start = end + 1
	end = start + strings.IndexRune(line[start:], ' ')
	host = line[start:end]

	start = end + 1
	end = start + strings.IndexRune(line[start:], ' ')
	if end > start {
		url = line[start:end]

		start = end + 1
		headers = parseHeaders(strings.TrimSuffix(strings.TrimPrefix(line[start:], "map["), "]"))
	} else {
		url = line[start:]
	}

	return timestamp.Unix(), []string{
		source,
		address,
		protocol,
		method,
		host,
		url,
	}, headers, nil
}

func parseHeaders(s string) map[string]string {
	headers := make(map[string]string)
	limit := len(s) - 1
	var (
		start, end int
		key, value string
	)
	for start < limit {
		end = start + strings.Index(s[start:], ":[")
		if end < start {
			break
		}
		key = s[start:end]

		start = end + 2
		end = start + strings.Index(s[start:], "] ")
		if end < start {
			value = s[start:limit]
			start = limit
		} else {
			value = s[start:end]
			start = end + 2
		}
		headers[key] = value
	}
	return headers
}
