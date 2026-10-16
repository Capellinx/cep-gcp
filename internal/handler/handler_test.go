package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cep-gcp/internal/viacep"
)

type fakeLoc struct {
	city string
	err  error
}

func (f fakeLoc) FindCity(context.Context, string) (string, error) { return f.city, f.err }

type fakeTemp struct {
	c   float64
	err error
}

func (f fakeTemp) CurrentCelsius(context.Context, string) (float64, error) { return f.c, f.err }

func do(t *testing.T, h http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestSuccess(t *testing.T) {
	h := New(fakeLoc{city: "São Paulo"}, fakeTemp{c: 28.5})
	rec := do(t, h, http.MethodGet, "/weather/01001000")

	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("content-type = %q", ct)
	}
	want := `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`
	if got := strings.TrimSpace(rec.Body.String()); got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}

func TestInvalidZipcode(t *testing.T) {
	h := New(fakeLoc{city: "X"}, fakeTemp{})
	for _, cep := range []string{"123", "123456789", "abcdefgh", "01001-000", "0100100a", "%2001001000"} {
		t.Run(cep, func(t *testing.T) {
			rec := do(t, h, http.MethodGet, "/weather/"+cep)
			if rec.Code != 422 {
				t.Errorf("status = %d, want 422", rec.Code)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != "invalid zipcode" {
				t.Errorf("body = %q", got)
			}
		})
	}
}

func TestZipcodeNotFound(t *testing.T) {
	h := New(fakeLoc{err: viacep.ErrNotFound}, fakeTemp{})
	rec := do(t, h, http.MethodGet, "/weather/99999999")
	if rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "can not find zipcode" {
		t.Errorf("body = %q", got)
	}
}

func TestUpstreamErrors(t *testing.T) {
	tests := map[string]struct {
		loc  fakeLoc
		temp fakeTemp
	}{
		"viacep falha":       {fakeLoc{err: errors.New("boom")}, fakeTemp{}},
		"weather falha":      {fakeLoc{city: "X"}, fakeTemp{err: errors.New("boom")}},
		"contexto cancelado": {fakeLoc{err: context.Canceled}, fakeTemp{}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			rec := do(t, New(tt.loc, tt.temp), http.MethodGet, "/weather/01001000")
			if rec.Code != 502 {
				t.Errorf("status = %d, want 502", rec.Code)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != "upstream error" {
				t.Errorf("body = %q", got)
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	h := New(fakeLoc{city: "X"}, fakeTemp{})
	rec := do(t, h, http.MethodPost, "/weather/01001000")
	if rec.Code != 405 {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}
