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
	"aletheiaware.com/netgo"
	"aletheiaware.com/netgo/handler"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/fs"
	"log"
	"math"
	"net/http"
	"net/url"
	"path"
	"strings"
)

//go:embed assets
var embeddedFS embed.FS

func Serve(name string) error {
	logFile, err := netgo.SetupLogging()
	if err != nil {
		return err
	}
	defer logFile.Close()
	log.Println("Log File:", logFile.Name())

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return err
	}

	// Create Multiplexer
	mux := http.NewServeMux()

	// Handle Static Assets
	staticFS, err := fs.Sub(embeddedFS, path.Join("assets", "static"))
	if err != nil {
		return err
	}
	handler.AttachStaticFSHandler(mux, staticFS, false, fmt.Sprintf("public, max-age=%d", 60*60*24*7*52)) // 52 week max-age

	// Parse Templates
	templateFS, err := fs.Sub(embeddedFS, path.Join("assets", "template"))
	if err != nil {
		return err
	}
	templates, err := template.ParseFS(templateFS, "*.go.html")
	if err != nil {
		return err
	}

	// TODO support exclude filters

	// TODO sessions - how a single address interacted with the server over time

	// Handle Request Data
	mux.Handle("/requests.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT tbl_requests.timestamp, tbl_requests.address, tbl_requests.protocol, tbl_requests.method, tbl_requests.host, tbl_requests.url FROM tbl_requests`
		raw += requestFiltersFromQuery(r.URL.Query())
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Requests{
			Start: math.MaxInt64,
		}
		for rows.Next() {
			r := &Request{}
			err = rows.Scan(&r.Timestamp, &r.Address, &r.Protocol, &r.Method, &r.Host, &r.URL)
			if err != nil {
				log.Fatal(err)
			}
			result.Total++
			if r.Timestamp < result.Start {
				result.Start = r.Timestamp
			}
			if r.Timestamp > result.End {
				result.End = r.Timestamp
			}
			result.Rows = append(result.Rows, r)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Fatal(err)
		}
	})))
	// Handle Address Data
	mux.Handle("/addresses.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT tbl_requests.address, COUNT(tbl_requests.id) AS count FROM tbl_requests`
		raw += requestFiltersFromQuery(r.URL.Query())
		raw += ` GROUP BY tbl_requests.address`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Addresses{}
		for rows.Next() {
			a := &Address{}
			err = rows.Scan(&a.Address, &a.Count)
			if err != nil {
				log.Fatal(err)
			}
			result.Total += a.Count
			if a.Count > result.Limit {
				result.Limit = a.Count
			}
			result.Rows = append(result.Rows, a)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Fatal(err)
		}
	})))
	// Handle Protocol Data
	mux.Handle("/protocols.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT tbl_requests.protocol, COUNT(tbl_requests.id) AS count FROM tbl_requests`
		raw += requestFiltersFromQuery(r.URL.Query())
		raw += ` GROUP BY tbl_requests.protocol`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Protocols{}
		for rows.Next() {
			p := &Protocol{}
			err = rows.Scan(&p.Protocol, &p.Count)
			if err != nil {
				log.Fatal(err)
			}
			result.Total += p.Count
			if p.Count > result.Limit {
				result.Limit = p.Count
			}
			result.Rows = append(result.Rows, p)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Fatal(err)
		}
	})))
	// Handle Method Data
	mux.Handle("/methods.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT tbl_requests.method, COUNT(tbl_requests.id) AS count FROM tbl_requests`
		raw += requestFiltersFromQuery(r.URL.Query())
		raw += ` GROUP BY tbl_requests.method`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Methods{}
		for rows.Next() {
			m := &Method{}
			err = rows.Scan(&m.Method, &m.Count)
			if err != nil {
				log.Fatal(err)
			}
			result.Total += m.Count
			if m.Count > result.Limit {
				result.Limit = m.Count
			}
			result.Rows = append(result.Rows, m)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Fatal(err)
		}
	})))
	// Handle URL Data
	mux.Handle("/urls.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT tbl_requests.url, COUNT(tbl_requests.id) AS count FROM tbl_requests`
		raw += requestFiltersFromQuery(r.URL.Query())
		raw += ` GROUP BY tbl_requests.url`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &URLs{}
		for rows.Next() {
			u := &URL{}
			err = rows.Scan(&u.URL, &u.Count)
			if err != nil {
				log.Fatal(err)
			}
			result.Total += u.Count
			if u.Count > result.Limit {
				result.Limit = u.Count
			}
			result.Rows = append(result.Rows, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Fatal(err)
		}
	})))
	// Handle Header Key Data
	mux.Handle("/header-keys.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT tbl_headers.key, COUNT(tbl_headers.id) AS count FROM tbl_headers`
		raw += headerFiltersFromQuery(r.URL.Query())
		raw += ` GROUP BY tbl_headers.key`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Headers{}
		for rows.Next() {
			h := &Header{}
			err = rows.Scan(&h.Key, &h.Count)
			if err != nil {
				log.Fatal(err)
			}
			result.Total += h.Count
			if h.Count > result.Limit {
				result.Limit = h.Count
			}
			result.Rows = append(result.Rows, h)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Fatal(err)
		}
	})))
	// Handle Header Value Data
	mux.Handle("/header-values.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT tbl_headers.value, COUNT(tbl_headers.id) AS count FROM tbl_headers`
		raw += headerFiltersFromQuery(r.URL.Query())
		raw += ` GROUP BY tbl_headers.value`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Headers{}
		for rows.Next() {
			h := &Header{}
			err = rows.Scan(&h.Value, &h.Count)
			if err != nil {
				log.Fatal(err)
			}
			result.Total += h.Count
			if h.Count > result.Limit {
				result.Limit = h.Count
			}
			result.Rows = append(result.Rows, h)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Fatal(err)
		}
	})))

	// Handle Index
	mux.Handle("/", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "index.go.html", &struct {
			Live bool
		}{
			Live: netgo.IsLive(),
		}); err != nil {
			log.Println(err)
		}
	})))

	// Serve HTTP Requests
	log.Println("HTTP Server Listening on :80")
	if err := http.ListenAndServe(":80", mux); err != nil {
		return err
	}
	return nil
}

func requestFiltersFromQuery(query url.Values) (result string) {
	var headerfilters []string
	hkey := netgo.QueryParameter(query, "header-key")
	if hkey != "" {
		headerfilters = append(headerfilters, fmt.Sprintf(`tbl_headers.key = '%s'`, hkey))
	}
	hvalue := netgo.QueryParameter(query, "header-value")
	if hvalue != "" {
		headerfilters = append(headerfilters, fmt.Sprintf(`tbl_headers.value = '%s'`, hvalue))
	}
	if len(headerfilters) > 0 {
		result += ` INNER JOIN tbl_headers ON tbl_requests.id = tbl_headers.request AND ` + strings.Join(headerfilters, ` AND `)
	}

	var requestfilters []string
	start := netgo.ParseInt(netgo.QueryParameter(query, "start"))
	if start != 0 {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.timestamp > %d`, start))
	}
	end := netgo.ParseInt(netgo.QueryParameter(query, "end"))
	if end != 0 {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.timestamp < %d`, end))
	}
	address := netgo.QueryParameter(query, "address")
	if address != "" {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.address = '%s'`, address))
	}
	protocol := netgo.QueryParameter(query, "protocol")
	if protocol != "" {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.protocol = '%s'`, protocol))
	}
	method := netgo.QueryParameter(query, "method")
	if method != "" {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.method = '%s'`, method))
	}
	url := netgo.QueryParameter(query, "url")
	if url != "" {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.url = '%s'`, url))
	}
	if len(requestfilters) > 0 {
		result += ` WHERE ` + strings.Join(requestfilters, ` AND `)
	}
	return
}

func headerFiltersFromQuery(query url.Values) (result string) {
	var requestfilters []string
	start := netgo.ParseInt(netgo.QueryParameter(query, "start"))
	if start != 0 {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.timestamp > %d`, start))
	}
	end := netgo.ParseInt(netgo.QueryParameter(query, "end"))
	if end != 0 {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.timestamp < %d`, end))
	}
	address := netgo.QueryParameter(query, "address")
	if address != "" {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.address = '%s'`, address))
	}
	protocol := netgo.QueryParameter(query, "protocol")
	if protocol != "" {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.protocol = '%s'`, protocol))
	}
	method := netgo.QueryParameter(query, "method")
	if method != "" {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.method = '%s'`, method))
	}
	url := netgo.QueryParameter(query, "url")
	if url != "" {
		requestfilters = append(requestfilters, fmt.Sprintf(`tbl_requests.url = '%s'`, url))
	}
	if len(requestfilters) > 0 {
		result += ` INNER JOIN tbl_requests ON tbl_requests.id = tbl_headers.request AND ` + strings.Join(requestfilters, ` AND `)
	}

	var headerfilters []string
	hkey := netgo.QueryParameter(query, "header-key")
	if hkey != "" {
		headerfilters = append(headerfilters, fmt.Sprintf(`tbl_headers.key = '%s'`, hkey))
	}
	hvalue := netgo.QueryParameter(query, "header-value")
	if hvalue != "" {
		headerfilters = append(headerfilters, fmt.Sprintf(`tbl_headers.value = '%s'`, hvalue))
	}
	if len(headerfilters) > 0 {
		result += ` WHERE ` + strings.Join(headerfilters, ` AND `)
	}
	return
}
