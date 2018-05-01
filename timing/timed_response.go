package timing

import "net/http"

type TimedResponse struct {
	http.ResponseWriter
	Timer *Timer
}

func (tr *TimedResponse) WriteHeader(code int) {
	tr.Timer.Stop()
	tr.ResponseWriter.WriteHeader(code)
}
