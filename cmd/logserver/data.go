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

type Requests struct {
	Total int        `json:"total"`
	Start int64      `json:"start"`
	End   int64      `json:"end"`
	Rows  []*Request `json:"rows"`
}

type Request struct {
	Timestamp int64             `json:"timestamp"`
	Address   string            `json:"address"`
	Protocol  string            `json:"protocol"`
	Method    string            `json:"method"`
	Host      string            `json:"host"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
}

type Addresses struct {
	Total int        `json:"total"`
	Limit int        `json:"limit"`
	Rows  []*Address `json:"rows"`
}

type Address struct {
	Address string `json:"address"`
	Count   int    `json:"count"`
}

type Protocols struct {
	Total int         `json:"total"`
	Limit int         `json:"limit"`
	Rows  []*Protocol `json:"rows"`
}

type Protocol struct {
	Protocol string `json:"protocol"`
	Count    int    `json:"count"`
}

type Methods struct {
	Total int       `json:"total"`
	Limit int       `json:"limit"`
	Rows  []*Method `json:"rows"`
}

type Method struct {
	Method string `json:"method"`
	Count  int    `json:"count"`
}

type URLs struct {
	Total int    `json:"total"`
	Limit int    `json:"limit"`
	Rows  []*URL `json:"rows"`
}

type URL struct {
	URL   string `json:"url"`
	Count int    `json:"count"`
}

type Headers struct {
	Total int       `json:"total"`
	Limit int       `json:"limit"`
	Rows  []*Header `json:"rows"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Count int    `json:"count"`
}
