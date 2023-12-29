package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.lguedes.ft/internal/assert"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)
	rs := rr.Result()
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "OK", string(body))
}

func TestSecureHeaders(t *testing.T) {
	// Setup:
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil) // Make a request
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	tests := []struct {
		name string
		want string
		got  string
	}{
		{
			name: "Content-Security-Policy",
			want: "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
			got:  rs.Header.Get("Content-Security-Policy"),
		},
		{
			name: "Referrer-Policy",
			want: "origin-when-cross-origin",
			got:  rs.Header.Get("Referrer-Policy"),
		},
		{
			name: "X-Content-Type-Options",
			want: "nosniff",
			got:  rs.Header.Get("X-Content-Type-Options"),
		},
		{
			name: "X-frame-Options",
			want: "deny",
			got:  rs.Header.Get("X-frame-Options"),
		},
		{
			name: "X-XSS-Protection",
			want: "0",
			got:  rs.Header.Get("X-XSS-Protection"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
	assert.Equal(t, "OK", string(body))
}
