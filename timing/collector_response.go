package timing

import "net/http"

type CollectorResponse struct {
	http.ResponseWriter
	Collector *Collector
}

func (r *CollectorResponse) WriteHeader(statusCode int) {
	r.Collector.WriteHeader(r)
	r.ResponseWriter.WriteHeader(statusCode)
}
