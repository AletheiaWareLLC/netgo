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
	Id        int    `json:"id"`
	File      string `json:"file"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
	Address   string `json:"address"`
	Protocol  string `json:"protocol"`
	Method    string `json:"method"`
	Host      string `json:"host"`
	URL       string `json:"url"`
	Cookie    string `json:"cookie"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"useragent"`
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

type Cookies struct {
	Total int       `json:"total"`
	Limit int       `json:"limit"`
	Rows  []*Cookie `json:"rows"`
}

type Cookie struct {
	Cookie string `json:"cookie"`
	Count  int    `json:"count"`
}

type Referrers struct {
	Total int         `json:"total"`
	Limit int         `json:"limit"`
	Rows  []*Referrer `json:"rows"`
}

type Referrer struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

type UserAgents struct {
	Total int          `json:"total"`
	Limit int          `json:"limit"`
	Rows  []*UserAgent `json:"rows"`
}

type UserAgent struct {
	UserAgent string `json:"useragent"`
	Count     int    `json:"count"`
}
