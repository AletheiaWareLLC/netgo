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

package netgo_test

import (
	"github.com/AletheiaWareLLC/netgo"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHTTPSRedirect(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/foo/bar", nil)
		response := httptest.NewRecorder()

		netgo.HTTPSRedirect(map[string]bool{
			"/foo/bar": true,
		})(response, request)

		if response.Code != http.StatusTemporaryRedirect {
			t.Errorf("Wrong response code; expected 300, got '%s'", http.StatusText(response.Code))
		}

		actual := response.Body.String()
		expected := "<a href=\"https:///foo/bar\">Temporary Redirect</a>.\n\n"

		if actual != expected {
			t.Errorf("Wrong response; expected '%s', got '%s'", expected, actual)
		}
	})
	t.Run("not allowed", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/foo/bar", nil)
		response := httptest.NewRecorder()

		netgo.HTTPSRedirect(map[string]bool{})(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("Wrong response code; expected 404, got '%s'", http.StatusText(response.Code))
		}

		actual := response.Body.String()
		expected := "404 page not found\n"

		if actual != expected {
			t.Errorf("Wrong response; expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestStaticHandler(t *testing.T) {
	// TODO
}

func assertQueryParameter(t *testing.T, query url.Values, key, expected string) {
	t.Helper()
	result := netgo.GetQueryParameter(query, "foo")
	if result != expected {
		t.Fatalf("Incorrect query parameter; expected '%s', got '%s'", expected, result)
	}
}

func TestGetQueryParameter(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		query := make(url.Values)
		expected := ""
		assertQueryParameter(t, query, "foo", expected)
	})
	t.Run("1", func(t *testing.T) {
		query := make(url.Values)
		query["foo"] = []string{
			"bar",
		}
		expected := "bar"
		assertQueryParameter(t, query, "foo", expected)
	})
	t.Run("2", func(t *testing.T) {
		query := make(url.Values)
		query["foo"] = []string{
			"bar",
			"baz",
		}
		expected := "bar"
		assertQueryParameter(t, query, "foo", expected)
	})
}
