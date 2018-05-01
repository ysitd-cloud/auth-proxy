package handler

import (
	"context"
	"net/http"
	"time"
)

const requestTimeout = 30 * time.Second
const timingKey = "timing"

func newContext(r *http.Request) (context.Context, context.CancelFunc) {
	ctx := r.Context()
	if _, deadline := ctx.Deadline(); !deadline {
		return context.WithTimeout(ctx, requestTimeout)
	} else {
		return context.WithCancel(ctx)
	}
}
