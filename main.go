package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ClickRequest struct {
	Count    int `json:"count"`
	IncCount int `json:"incCount"`
	DecCount int `json:"decCount"`
}

func newMux(buttonValue prometheus.Gauge) *http.ServeMux {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		buttonValue,
	)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "static/index.html")
	})

	mux.HandleFunc("/click", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var click ClickRequest

		if err := json.NewDecoder(r.Body).Decode(&click); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		buttonValue.Set(float64(click.Count))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Button clicked!"))

		fmt.Println("Response got")
	})

	return mux
}

func main() {
	buttonValue := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "button_value",
		Help: "Current value of the button",
	})

	mux := newMux(buttonValue)

	fmt.Println("server listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
