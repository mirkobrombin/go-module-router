//go:build !ping
// +build !ping

package ping

import (
	"encoding/json"
	"strconv"

	"github.com/valyala/fasthttp"
)

type PingHandler struct{ svc PingService }

type pingResponse struct {
	Message []string `json:"message"`
}

func (h *PingHandler) Pong(ctx *fasthttp.RequestCtx) {
	n, err := strconv.Atoi(string(ctx.QueryArgs().Peek("times")))
	if err != nil || n < 1 {
		n = 1
	}

	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = h.svc.Pong()
	}

	payload, _ := json.Marshal(pingResponse{Message: out})

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	_, _ = ctx.Write(payload)
}
