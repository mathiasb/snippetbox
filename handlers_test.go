package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	rs := rr.Result()
	defer rs.Body.Close()

	checkPing(rs.StatusCode, rs, t)
}

func checkPing(code int, rs *http.Response, t *testing.T) {
	t.Helper()
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestPingE2E(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippetE2E(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running and end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	// Log the CSRF token value in our test output. To see the output from the
	// t.Log() command you need to run `go test` with the -v (verbose) flag
	// enabled.
	t.Log(csrfToken)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters)")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			if code != tt.wantCode {
				t.Errorf("want code %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestHTMXRequests(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Test that HTMX requests are handled correctly
	t.Run("HTMX Home Request", func(t *testing.T) {
		// Make an HTMX request to the home page
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("HX-Request", "true")
		
		rs, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer rs.Body.Close()
		
		if rs.StatusCode != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
		}
	})
	
	t.Run("HTMX Create Snippet Request", func(t *testing.T) {
		// Make an HTMX request to create a snippet
		form := url.Values{}
		form.Add("title", "Test Title")
		form.Add("content", "Test content")
		form.Add("expires", "365")
		form.Add("csrf_token", "test-token")
		
		req, err := http.NewRequest("POST", "/snippet/create", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("HX-Request", "true")
		
		rs, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer rs.Body.Close()
		
		// For HTMX requests, we expect a 200 OK with HX-Redirect header
		if rs.StatusCode != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
		}
		
		// Check that the response has the correct redirect header for HTMX
		redirect := rs.Header.Get("HX-Redirect")
		if redirect == "" {
			t.Error("expected HX-Redirect header for HTMX request")
		}
	})
	
	t.Run("HTMX Login Request", func(t *testing.T) {
		// Make an HTMX request to login
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", "test-token")
		
		req, err := http.NewRequest("POST", "/user/login", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("HX-Request", "true")
		
		rs, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer rs.Body.Close()
		
		// For HTMX requests, we expect a 200 OK with HX-Redirect header
		if rs.StatusCode != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
		}
		
		// Check that the response has the correct redirect header for HTMX
		redirect := rs.Header.Get("HX-Redirect")
		if redirect == "" {
			t.Error("expected HX-Redirect header for HTMX request")
		}
	})
}
```