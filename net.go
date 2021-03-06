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

package netgo

import (
	"log"
	"net/http"
	"net/url"
	"path"
)

const HTTPS = "HTTPS"

func HTTPSRedirect(host string, paths map[string]bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		allowed, ok := paths[r.URL.Path]
		if allowed && ok && r.Host == host {
			target := "https://" + r.Host + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			log.Println(r.RemoteAddr, r.Proto, r.Method, r.Host, r.URL.Path, r.Header, "redirected to", target)
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		} else {
			log.Println(r.RemoteAddr, r.Proto, r.Method, r.Host, r.URL.Path, r.Header, "not found")
			http.NotFound(w, r)
		}
	}
}

func StaticHandler(directory string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Proto, r.Method, r.Host, r.URL.Path, r.Header)
		switch r.Method {
		case "HEAD":
			fallthrough
		case "GET":
			http.ServeFile(w, r, path.Join(directory, r.URL.Path))
		default:
			log.Println("Unsupported method", r.Method)
		}
	}
}

func GetQueryParameter(query url.Values, parameter string) string {
	if results, ok := query[parameter]; ok && len(results) > 0 {
		return results[0]
	}
	return ""
}
