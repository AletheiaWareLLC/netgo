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

func Serve() error {
	logFile, err := netgo.SetupLogging()
	if err != nil {
		return err
	}
	defer logFile.Close()
	log.Println("Log File:", logFile.Name())

	db, err := sql.Open("sqlite3", "./log.db")
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
		raw := `SELECT * FROM tbl_requests`
		if filters := filtersFromQuery(r.URL.Query()); len(filters) > 0 {
			raw += ` WHERE ` + strings.Join(filters, ` AND `)
		}
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
			err = rows.Scan(&r.Id, &r.File, &r.Timestamp, &r.Source, &r.Address, &r.Protocol, &r.Method, &r.Host, &r.URL, &r.Cookie, &r.Referrer, &r.UserAgent)
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
		raw := `SELECT address, COUNT(id) AS count FROM tbl_requests`
		if filters := filtersFromQuery(r.URL.Query()); len(filters) > 0 {
			raw += ` WHERE ` + strings.Join(filters, ` AND `)
		}
		raw += ` GROUP BY address`
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
		raw := `SELECT protocol, COUNT(id) AS count FROM tbl_requests`
		if filters := filtersFromQuery(r.URL.Query()); len(filters) > 0 {
			raw += ` WHERE ` + strings.Join(filters, ` AND `)
		}
		raw += ` GROUP BY protocol`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Protocols{}
		for rows.Next() {
			a := &Protocol{}
			err = rows.Scan(&a.Protocol, &a.Count)
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
	// Handle Method Data
	mux.Handle("/methods.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT method, COUNT(id) AS count FROM tbl_requests`
		if filters := filtersFromQuery(r.URL.Query()); len(filters) > 0 {
			raw += ` WHERE ` + strings.Join(filters, ` AND `)
		}
		raw += ` GROUP BY method`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Methods{}
		for rows.Next() {
			a := &Method{}
			err = rows.Scan(&a.Method, &a.Count)
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
	// Handle URL Data
	mux.Handle("/urls.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT url, COUNT(id) AS count FROM tbl_requests`
		if filters := filtersFromQuery(r.URL.Query()); len(filters) > 0 {
			raw += ` WHERE ` + strings.Join(filters, ` AND `)
		}
		raw += ` GROUP BY url`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &URLs{}
		for rows.Next() {
			a := &URL{}
			err = rows.Scan(&a.URL, &a.Count)
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
	// Handle Cookie Data
	mux.Handle("/cookies.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT cookie, COUNT(id) AS count FROM tbl_requests`
		if filters := filtersFromQuery(r.URL.Query()); len(filters) > 0 {
			raw += ` WHERE ` + strings.Join(filters, ` AND `)
		}
		raw += ` GROUP BY cookie`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Cookies{}
		for rows.Next() {
			a := &Cookie{}
			err = rows.Scan(&a.Cookie, &a.Count)
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
	// Handle Referrer Data
	mux.Handle("/referrers.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT referrer, COUNT(id) AS count FROM tbl_requests`
		if filters := filtersFromQuery(r.URL.Query()); len(filters) > 0 {
			raw += ` WHERE ` + strings.Join(filters, ` AND `)
		}
		raw += ` GROUP BY referrer`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &Referrers{}
		for rows.Next() {
			a := &Referrer{}
			err = rows.Scan(&a.Referrer, &a.Count)
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
	// Handle UserAgent Data
	mux.Handle("/useragents.json", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `SELECT useragent, COUNT(id) AS count FROM tbl_requests`
		if filters := filtersFromQuery(r.URL.Query()); len(filters) > 0 {
			raw += ` WHERE ` + strings.Join(filters, ` AND `)
		}
		raw += ` GROUP BY useragent`
		raw += ` ORDER BY count DESC`
		raw += ` LIMIT 1000`
		rows, err := db.Query(raw)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		result := &UserAgents{}
		for rows.Next() {
			a := &UserAgent{}
			err = rows.Scan(&a.UserAgent, &a.Count)
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

func filtersFromQuery(query url.Values) (filters []string) {
	start := netgo.ParseInt(netgo.QueryParameter(query, "start"))
	if start != 0 {
		filters = append(filters, fmt.Sprintf(`timestamp > %d`, start))
	}
	end := netgo.ParseInt(netgo.QueryParameter(query, "end"))
	if end != 0 {
		filters = append(filters, fmt.Sprintf(`timestamp < %d`, end))
	}
	address := netgo.QueryParameter(query, "address")
	if address != "" {
		filters = append(filters, fmt.Sprintf(`address = '%s'`, address))
	}
	protocol := netgo.QueryParameter(query, "protocol")
	if protocol != "" {
		filters = append(filters, fmt.Sprintf(`protocol = '%s'`, protocol))
	}
	method := netgo.QueryParameter(query, "method")
	if method != "" {
		filters = append(filters, fmt.Sprintf(`method = '%s'`, method))
	}
	cookie := netgo.QueryParameter(query, "cookie")
	if cookie != "" {
		filters = append(filters, fmt.Sprintf(`cookie = '%s'`, cookie))
	}
	url := netgo.QueryParameter(query, "url")
	if url != "" {
		filters = append(filters, fmt.Sprintf(`url = '%s'`, url))
	}
	referrer := netgo.QueryParameter(query, "referrer")
	if referrer != "" {
		filters = append(filters, fmt.Sprintf(`referrer = '%s'`, referrer))
	}
	useragent := netgo.QueryParameter(query, "useragent")
	if useragent != "" {
		filters = append(filters, fmt.Sprintf(`useragent = '%s'`, useragent))
	}
	return
}
