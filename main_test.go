package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func newTestServer() (*http.ServeMux, prometheus.Gauge) {
	buttonValue := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "button_value",
		Help: "Current value of the button",
	})

	return newMux(buttonValue), buttonValue
}

func TestClickSetsButtonValue(t *testing.T) {
	mux, buttonValue := newTestServer()

	payload := ClickRequest{
		Count:    7,
		IncCount: 7,
		DecCount: 2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/click", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s",
			http.StatusOK, rec.Code, rec.Body.String())
	}

	if got := rec.Body.String(); got != "Button clicked!" {
		t.Fatalf("expected response body %q, got %q",
			"Button clicked!", got)
	}

	expectedMetric := `
# HELP button_value Current value of the button
# TYPE button_value gauge
button_value 7
`

	if err := testutil.CollectAndCompare(
		buttonValue,
		strings.NewReader(expectedMetric),
		"button_value",
	); err != nil {
		t.Fatalf("unexpected metric value: %v", err)
	}
}

func TestClickRejectsInvalidJSON(t *testing.T) {
	mux, _ := newTestServer()

	req := httptest.NewRequest(
		http.MethodPost,
		"/click",
		strings.NewReader(`{"count":`),
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d; body: %s",
			http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestClickRejectsGET(t *testing.T) {
	mux, _ := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/click", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d; body: %s",
			http.StatusMethodNotAllowed, rec.Code, rec.Body.String())
	}

	if got := rec.Header().Get("Allow"); got != http.MethodPost {
		t.Fatalf("expected Allow header %q, got %q", http.MethodPost, got)
	}
}

func TestMetricsContainsButtonValue(t *testing.T) {
	mux, _ := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "button_value") {
		t.Fatalf("expected /metrics response to contain button_value")
	}
}

func TestHomePage(t *testing.T) {
	mux, _ := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s",
			http.StatusOK, rec.Code, rec.Body.String())
	}

	if !strings.Contains(rec.Body.String(), "Simple page") {
		t.Fatalf("expected page content to contain title %q",
			"Simple page")
	}
}

func TestUnknownRouteReturnsNotFound(t *testing.T) {
	mux, _ := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
