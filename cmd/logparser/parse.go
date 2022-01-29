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
	"runtime"
	"strings"
)

const (
	CREATE_QUERY = `CREATE TABLE IF NOT EXISTS tbl_requests (
    id INTEGER NOT NULL PRIMARY KEY,
    file TEXT,
    timestamp INT UNSIGNED NOT NULL,
    source TEXT,
    address TEXT,
    protocol TEXT,
    method TEXT,
    host TEXT,
    url TEXT,
    cookie TEXT,
    referrer TEXT,
    useragent TEXT
);`
	INSERT_QUERY = `INSERT INTO tbl_requests
(file, timestamp, source, address, protocol, method, host, url, cookie, referrer, useragent)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
)

func Parse(name string, sources, dirs []string) (int, error) {
	// Create sqlite database
	db, err := openDatabase(name)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Limit to one connection to avoid 'database is locked' error
	db.SetMaxOpenConns(1)

	jobs := make(chan string, 10000)
	defer close(jobs)

	results := make(chan int, 10000)
	defer close(results)

	ignores := make(chan string, 10000)
	defer close(ignores)

	go func() {
		os.Remove(".ignored")
		// Create file of ignored logs
		ignored, err := os.Create(".ignored")
		if err != nil {
			log.Fatal(err)
		}
		defer ignored.Close()

		for i := range ignores {
			if _, err := ignored.WriteString(i + "\n"); err != nil {
				log.Fatal(err)
			}
		}

		if err := ignored.Sync(); err != nil {
			log.Fatal(err)
		}
	}()

	work := func() {
		for j := range jobs {
			count, err := parseLog(db, sources, j, ignores)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(j, count)
			results <- count
		}
	}
	// Spawns a worker thread for each available process.
	for w := 0; w < runtime.GOMAXPROCS(0); w++ {
		go work()
	}

	var count int
	for _, dir := range dirs {
		// Scan directory for logs
		ls, err := os.ReadDir(dir)
		if err != nil {
			return 0, err
		}

		// Parse each log file in a separate worker
		for _, l := range ls {
			name := l.Name()
			jobs <- path.Join(dir, name)
		}

		// Wait for all workers to complete
		for i := 0; i < len(ls); i++ {
			count += <-results
		}
	}

	return count, nil
}

func parseLog(db *sql.DB, sources []string, name string, ignores chan string) (int, error) {
	var count int

	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
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
			if _, err := db.Exec(INSERT_QUERY, name, timestamp, request[0], request[1], request[2], request[3], request[4], request[5], headers["Cookie"], headers["Referer"], headers["User-Agent"]); err != nil {
				return 0, err
			}
			count++
		} else {
			ignores <- line
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
	// Create table for requests
	if _, err = db.Exec(CREATE_QUERY); err != nil {
		return nil, err
	}
	return db, nil
}
