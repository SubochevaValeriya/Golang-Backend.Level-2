package main

import (
	"database/sql"
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

	router.HandleFunc("/identity", measurable(DBMethods)).Methods(http.MethodGet)
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
func DBMethods(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	sql := "INSERT INTO cities(name, population) VALUES ('Moscow', 12506000)"
	_, err = db.Exec(sql)
	_, err2 := db.Query("SELECT * FROM cities WHERE id = ?", 1)

	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
