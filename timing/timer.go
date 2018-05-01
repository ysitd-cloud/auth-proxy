package timing

import (
	"fmt"
	"net/http"
	"time"
)

type Timer struct {
	Begin    time.Time
	End      time.Time
	Ended    bool
	Key      string
	Desc     string
	Children []*Timer
}

func NewTimer(key, desc string) *Timer {
	return &Timer{
		Key:      key,
		Desc:     desc,
		Ended:    false,
		Children: make([]*Timer, 0),
	}
}

func (t *Timer) Start() {
	t.Begin = time.Now()
}

func (t *Timer) Stop() {
	if !t.Ended {
		t.Ended = true
		t.End = time.Now()
	}
}

func (t *Timer) AddChild(child *Timer) {
	t.Children = append(t.Children, child)
}

func (t *Timer) WriteHeader(w http.ResponseWriter) {
	t.Stop()
	for _, child := range t.Children {
		child.WriteHeader(w)
	}
	duration := t.End.Sub(t.Begin)
	dur := float64(duration/time.Millisecond) + float64(duration%time.Millisecond)/float64(time.Millisecond)
	header := fmt.Sprintf("%s;desc=\"%s\";dur=%f", t.Key, t.Desc, dur)
	w.Header().Add("Server-Timing", header)
}

func (t *Timer) Response(w http.ResponseWriter) *TimedResponse {
	return &TimedResponse{
		ResponseWriter: w,
		Timer:          t,
	}
}
