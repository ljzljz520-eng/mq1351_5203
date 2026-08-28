package web

import (
	"fmt"
	"net/http"
	"strings"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(data)
	r.bytes += n
	return n, err
}

func withResponseMetrics(next http.Handler, observe func(status, bytes int)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		observe(recorder.status, recorder.bytes)
	})
}

func acceptsJSON(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Accept")), "application/json")
}

func clientLabel(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("X-Visitor")); value != "" {
		return value
	}
	return "匿名琴友"
}

func statusText(status int) string {
	if status >= 200 && status < 300 {
		return "success"
	}
	if status >= 400 {
		return "error"
	}
	return "other"
}

func formatResponseMetric(status, bytes int) string {
	return fmt.Sprintf("%s:%d", statusText(status), bytes)
}
