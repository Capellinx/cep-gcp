package viacep

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, status int, body string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json/" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return &Client{BaseURL: srv.URL, HTTP: srv.Client()}
}

func TestFindCity_Success(t *testing.T) {
	c := newTestClient(t, 200, `{"cep":"01001-000","localidade":"São Paulo","uf":"SP"}`)
	got, err := c.FindCity(context.Background(), "01001000")
	if err != nil {
		t.Fatal(err)
	}
	if got != "São Paulo" {
		t.Errorf("city = %q", got)
	}
}

func TestFindCity_NotFound(t *testing.T) {
	for name, body := range map[string]string{
		"erro booleano": `{"erro": true}`,
		"erro string":   `{"erro": "true"}`,
	} {
		t.Run(name, func(t *testing.T) {
			c := newTestClient(t, 200, body)
			_, err := c.FindCity(context.Background(), "01001000")
			if !errors.Is(err, ErrNotFound) {
				t.Errorf("err = %v, want ErrNotFound", err)
			}
		})
	}
}

func TestFindCity_UpstreamErrors(t *testing.T) {
	tests := map[string]struct {
		status int
		body   string
	}{
		"5xx":              {500, `oops`},
		"json invalido":    {200, `not json`},
		"sem localidade":   {200, `{"cep":"01001-000"}`},
		"localidade vazia": {200, `{"localidade":""}`},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := newTestClient(t, tt.status, tt.body)
			_, err := c.FindCity(context.Background(), "01001000")
			if err == nil || errors.Is(err, ErrNotFound) {
				t.Errorf("err = %v, want non-NotFound error", err)
			}
		})
	}
}
