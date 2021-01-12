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

package main

import (
	"aletheiaware.com/netgo"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "start":
			if err := start(); err != nil {
				log.Println(err)
				return
			}
		default:
			log.Println("Cannot handle", os.Args[1])
		}
	} else {
		PrintUsage(os.Stdout)
	}
}

func start() error {
	store, ok := os.LookupEnv("LOG_DIRECTORY")
	if !ok {
		store = "logs"
	}
	if err := os.MkdirAll(store, os.ModePerm); err != nil {
		return err
	}
	logFile, err := os.OpenFile(path.Join(store, time.Now().UTC().Format(time.RFC3339)), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Log File:", logFile.Name())

	// Serve Web Requests
	mux := http.NewServeMux()
	mux.HandleFunc("/", netgo.StaticHandler("html/static"))

	if https, ok := os.LookupEnv(netgo.HTTPS); ok && https == "true" {
		certificates, ok := os.LookupEnv("CERTIFICATE_DIRECTORY")
		if !ok {
			certificates = "certificates"
		}
		log.Println("Certificate Directory:", certificates)

		host, ok := os.LookupEnv("HOST")
		if !ok {
			return errors.New("Missing HOST environment variable")
		}

		routeMap := make(map[string]bool)

		routes, ok := os.LookupEnv("ROUTES")
		if ok {
			for _, route := range strings.Split(routes, ",") {
				routeMap[route] = true
			}
		}

		// Redirect HTTP Requests to HTTPS
		go http.ListenAndServe(":80", http.HandlerFunc(netgo.HTTPSRedirect(host, routeMap)))

		// Serve HTTPS Requests
		config := &tls.Config{MinVersion: tls.VersionTLS10}
		server := &http.Server{Addr: ":443", Handler: mux, TLSConfig: config}
		return server.ListenAndServeTLS(path.Join(certificates, "fullchain.pem"), path.Join(certificates, "privkey.pem"))
	} else {
		log.Println("HTTP Server Listening on :80")
		return http.ListenAndServe(":80", mux)
	}
}

func PrintUsage(output io.Writer) {
	fmt.Fprintln(output, "Net Server Usage:")
	fmt.Fprintln(output, "\tserver - display usage")
	fmt.Fprintln(output, "\tserver start - starts the server")
}
