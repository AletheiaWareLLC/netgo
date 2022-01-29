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
	"aletheiaware.com/netgo/handler"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
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
	logFile, err := netgo.SetupLogging()
	if err != nil {
		return err
	}
	defer logFile.Close()
	log.Println("Log File:", logFile.Name())

	content, ok := os.LookupEnv("CONTENT_DIRECTORY")
	if !ok {
		content = "html/static"
	}

	// Serve Web Requests
	mux := http.NewServeMux()
	mux.Handle("/", handler.Log(handler.StaticDir(content, true)))

	if netgo.IsSecure() {
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
		config := &tls.Config{MinVersion: tls.VersionTLS12}
		server := &http.Server{Addr: ":443", Handler: mux, TLSConfig: config}
		return server.ListenAndServeTLS(filepath.Join(certificates, "fullchain.pem"), filepath.Join(certificates, "privkey.pem"))
	} else {
		log.Println("HTTP Server Listening on :80")
		return http.ListenAndServe(":80", mux)
	}
}

func PrintUsage(output io.Writer) {
	fmt.Fprintln(output, "Net Server Usage:")
	fmt.Fprintln(output, "\tnetserver - display usage")
	fmt.Fprintln(output, "\tnetserver start - starts the server")
}
