package middleware

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"sync"
	"time"
)

type IPLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	limit    int
	duration time.Duration
}

type visitor struct {
	requests int
	lastSeen time.Time
}

func NewIPLimiter(limit int, duration time.Duration) *IPLimiter {
	return &IPLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		duration: duration,
	}
}

func (i *IPLimiter) isLimited(ip string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.visitors[ip]
	if !exists {
		i.visitors[ip] = &visitor{
			requests: 1,
			lastSeen: time.Now(),
		}
		return false
	}

	// Check if duration has passed
	if time.Since(v.lastSeen) > i.duration {
		v.requests = 1
		v.lastSeen = time.Now()
		return false
	}

	// Increment request count
	v.requests++
	if v.requests > i.limit {
		return true
	}

	return false
}

func (i *IPLimiter) CleanUp() {
	for {
		time.Sleep(i.duration)
		i.mu.Lock()
		for ip, v := range i.visitors {
			if time.Since(v.lastSeen) > i.duration {
				delete(i.visitors, ip)
			}
		}
		i.mu.Unlock()
	}
}

func MiddlewareIPLimiter(ipLimiter *IPLimiter) beego.FilterFunc {
	return func(ctx *context.Context) {
		ip := ctx.Input.IP()
		if ipLimiter.isLimited(ip) {
			ctx.ResponseWriter.WriteHeader(429) // HTTP 429 Too Many Requests
			ctx.Output.JSON(map[string]interface{}{
				"statusError": "Too many request",
				"wait":        ipLimiter.duration,
			}, true, true)
			return
		}
	}
}
