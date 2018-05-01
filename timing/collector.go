package timing

import "net/http"

type Collector struct {
	Timers []*Timer
}

func NewCollector() *Collector {
	return &Collector{
		Timers: make([]*Timer, 0),
	}
}

func (c *Collector) New(key, desc string) (timer *Timer) {
	timer = NewTimer(key, desc)
	c.Add(timer)
	return
}

func (c *Collector) Add(timer *Timer) {
	c.Timers = append(c.Timers, timer)
}

func (c *Collector) WriteHeader(w http.ResponseWriter) {
	for _, timer := range c.Timers {
		timer.WriteHeader(w)
	}
}

func (c *Collector) Response(w http.ResponseWriter) *CollectorResponse {
	return &CollectorResponse{
		ResponseWriter: w,
		Collector:      c,
	}
}
