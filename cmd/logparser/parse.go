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
	"bufio"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path"
	"strings"
)

const (
	CREATE_FILES_QUERY = `CREATE TABLE IF NOT EXISTS tbl_files (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);`
	CREATE_REQUESTS_QUERY = `CREATE TABLE IF NOT EXISTS tbl_requests (
    id INTEGER NOT NULL PRIMARY KEY,
    file INT NULL,
    timestamp INT UNSIGNED NOT NULL,
    source TEXT,
    address TEXT,
    protocol TEXT,
    method TEXT,
    host TEXT,
    url TEXT,
    FOREIGN KEY (file) REFERENCES tbl_files(id)
);`
	CREATE_HEADERS_QUERY = `CREATE TABLE IF NOT EXISTS tbl_headers (
	id INTEGER NOT NULL PRIMARY KEY,
	request INT NULL,
	key TEXT,
	value TEXT,
	FOREIGN KEY (request) REFERENCES tbl_requests(id)
);`
	SELECT_FILE_QUERY = `SELECT * FROM tbl_files WHERE name = ?;`
	INSERT_FILE_QUERY = `INSERT INTO tbl_files
(name)
VALUES
(?);`
	INSERT_REQUEST_QUERY = `INSERT INTO tbl_requests
(file, timestamp, source, address, protocol, method, host, url)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?);`
	INSERT_HEADER_QUERY = `INSERT INTO tbl_headers
(request, key, value)
VALUES
(?, ?, ?);`
)

func Parse(name string, sources, dirs []string) (int, error) {
	// Create sqlite database
	db, err := openDatabase(name)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	os.Remove(".ignored")
	// Create file of ignored logs
	ignored, err := os.Create(".ignored")
	if err != nil {
		log.Fatal(err)
	}
	defer ignored.Close()

	var count int
	for _, dir := range dirs {
		// Scan directory for logs
		ls, err := os.ReadDir(dir)
		if err != nil {
			return 0, err
		}

		for _, l := range ls {
			name := path.Join(dir, l.Name())
			// Check if file has already been parsed
			row := db.QueryRow(SELECT_FILE_QUERY, name)
			var (
				id int
				n  string
			)
			err := row.Scan(&id, &n)
			if err == nil {
				// File already parsed
				continue
			} else if err != sql.ErrNoRows {
				return 0, err
			}
			c, err := parseLog(db, sources, name, func(l string) {
				if _, err := ignored.WriteString(l + "\n"); err != nil {
					log.Fatal(err)
				}
			})
			if err != nil {
				return 0, err
			}
			log.Println(name, c)
			count += c
		}
	}

	if err := ignored.Sync(); err != nil {
		log.Fatal(err)
	}

	return count, nil
}

func parseLog(db *sql.DB, sources []string, name string, onIgnore func(string)) (int, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	result, err := db.Exec(INSERT_FILE_QUERY, name)
	if err != nil {
		return 0, err
	}
	fileId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	var count int
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if netgo.IsRequestLog(sources, line) {
			timestamp, request, headers, err := netgo.ParseRequestLog(line)
			if err != nil {
				return 0, err
			}
			result, err := db.Exec(INSERT_REQUEST_QUERY, fileId, timestamp, request[0], request[1], request[2], request[3], request[4], request[5])
			if err != nil {
				return 0, err
			}
			requestId, err := result.LastInsertId()
			if err != nil {
				return 0, err
			}
			for k, v := range headers {
				if _, err := db.Exec(INSERT_HEADER_QUERY, requestId, k, v); err != nil {
					return 0, err
				}
			}
			count++
		} else {
			onIgnore(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

func openDatabase(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	// Create table for files
	if _, err = db.Exec(CREATE_FILES_QUERY); err != nil {
		return nil, err
	}
	// Create table for requests
	if _, err = db.Exec(CREATE_REQUESTS_QUERY); err != nil {
		return nil, err
	}
	// Create table for headers
	if _, err = db.Exec(CREATE_HEADERS_QUERY); err != nil {
		return nil, err
	}
	return db, nil
}
