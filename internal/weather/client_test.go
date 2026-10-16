package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testKey = "secret-key-123"

func TestCurrentCelsius_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/current.json" {
			t.Errorf("path = %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("key") != testKey {
			t.Errorf("key = %q", q.Get("key"))
		}
		if q.Get("q") != "São Paulo" { // acento e espaço devem chegar íntegros
			t.Errorf("q = %q", q.Get("q"))
		}
		_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIKey: testKey, HTTP: srv.Client()}
	got, err := c.CurrentCelsius(context.Background(), "São Paulo")
	if err != nil {
		t.Fatal(err)
	}
	if got != 28.5 {
		t.Errorf("temp = %v", got)
	}
}

func TestCurrentCelsius_NegativeAndZero(t *testing.T) {
	for _, body := range []string{`{"current":{"temp_c":-3.2}}`, `{"current":{"temp_c":0}}`} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(body))
		}))
		c := &Client{BaseURL: srv.URL, APIKey: testKey, HTTP: srv.Client()}
		if _, err := c.CurrentCelsius(context.Background(), "X"); err != nil {
			t.Errorf("body %s: %v", body, err)
		}
		srv.Close()
	}
}

func TestCurrentCelsius_Errors(t *testing.T) {
	tests := map[string]struct {
		status int
		body   string
	}{
		"4xx":           {400, `{"error":{"code":1006}}`},
		"5xx":           {500, `oops`},
		"json invalido": {200, `nope`},
		"sem current":   {200, `{}`},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()
			c := &Client{BaseURL: srv.URL, APIKey: testKey, HTTP: srv.Client()}
			if _, err := c.CurrentCelsius(context.Background(), "X"); err == nil {
				t.Error("want error")
			}
		})
	}
}

func TestCurrentCelsius_NetworkErrorDoesNotLeakKey(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close() // porta fechada -> erro de conexão

	c := &Client{BaseURL: url, APIKey: testKey, HTTP: &http.Client{}}
	_, err := c.CurrentCelsius(context.Background(), "X")
	if err == nil {
		t.Fatal("want error")
	}
	if strings.Contains(err.Error(), testKey) {
		t.Errorf("error leaks api key: %v", err)
	}
}
