package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractBearerToken(t *testing.T) {
	cases := []struct {
		header string
		want   string
		ok     bool
	}{
		{"", "", false},
		{"Bearer abc123", "abc123", true},
		{"bearer abc123", "abc123", true},
		{"  Bearer    spaced-token  ", "", false}, // leading whitespace breaks prefix match
		{"Bearer ", "", true},                     // present-but-empty token is still bearer-shaped (treated as invalid by repo)
		{"Basic Zm9vOmJhcg==", "", false},
	}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodGet, "/web/api/stats", nil)
		if c.header != "" {
			r.Header.Set("Authorization", c.header)
		}
		got, ok := extractBearerToken(r)
		if ok != c.ok || got != c.want {
			t.Errorf("extractBearerToken(%q) = (%q,%v); want (%q,%v)", c.header, got, ok, c.want, c.ok)
		}
	}
}

func TestIsAPIRequest(t *testing.T) {
	cases := []struct {
		path   string
		accept string
		want   bool
	}{
		{"/web/api/stats", "", true},
		{"/web/dashboard", "", false},
		{"/web/dashboard", "application/json", true},
		{"/web/dashboard", "text/html", false},
		{"/web/login", "", false},
	}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodGet, c.path, nil)
		if c.accept != "" {
			r.Header.Set("Accept", c.accept)
		}
		if got := isAPIRequest(r); got != c.want {
			t.Errorf("isAPIRequest(path=%q, accept=%q) = %v; want %v", c.path, c.accept, got, c.want)
		}
	}
}
