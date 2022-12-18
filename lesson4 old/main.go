package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

var duration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name:       "duration_seconds",
	Help:       "Summary of request duration in seconds",
	Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
},
	[]string{"labelHandler", "labelMethod", "labelStatus"})
var errorsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "errors_total",
	Help: "Total number of errors"},
	[]string{"labelHandler", "labelMethod", "labelStatus"})
var requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "request_total",
	Help: "Total number of requests"},
	[]string{"labelHandler", "labelMethod"})

var MeasurableHandler = func(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		m := r.Method
		p := r.URL.Path
		requestsTotal.WithLabelValues(p, m).Inc()
		mw := NewResponseWriter(w)
		h(mw, r)
		if mw.Status/100 > 3 {
			errorsTotal.WithLabelValues(p, m, strconv.Itoa(mw.Status)).Inc()
		}
		duration.WithLabelValues(p, m, strconv.Itoa(mw.Status)).Observe(time.Since(t).Seconds())
	}
}

type responseWriter struct {
	http.ResponseWriter
	Status int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

var (
	measurable = MeasurableHandler
	router     = mux.NewRouter()

	web = http.Server{
		Handler: router,
	}
)

func init() {
	router.HandleFunc("/identity", measurable(GetIdentityHandler)).Methods(http.MethodGet)
}
func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}
func GetIdentityHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") == "admin_secret_token" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}
