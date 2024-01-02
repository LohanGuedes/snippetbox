package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.lguedes.ft/internal/assert"
)

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/signup")
	t.Log(body)
	csrfToken := extractCSRFToken(t, body)
	t.Logf("CSRF token in: %q", csrfToken)
}

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	status, _, body := ts.Get(t, "/ping")

	assert.Equal(t, "OK", body)
	assert.Equal(t, http.StatusOK, status)
}

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
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

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantBody string
		wantCode int
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An Old Silent Pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/999",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.Get(t, tt.urlPath)
			assert.Equal(t, tt.wantCode, code)

			if tt.wantBody != "" {
				assert.StringContains(t, tt.wantBody, body)
			}
		})
	}
}
