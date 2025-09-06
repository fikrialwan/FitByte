package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/fikrialwan/FitByte/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	return i
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter

	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()
	return limiter
}

func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Global rate limiter
var GlobalRateLimiter *IPRateLimiter

// InitGlobalRateLimiter initializes the global rate limiter with config
func InitGlobalRateLimiter(cfg *config.Config) {
	rateLimit := rate.Every(time.Microsecond)
	if cfg.RateLimitPerSecond > 0 {
		rateLimit = rate.Limit(cfg.RateLimitPerSecond)
	}

	burst := 100
	if cfg.RateLimitBurst > 0 {
		burst = cfg.RateLimitBurst
	}

	GlobalRateLimiter = NewIPRateLimiter(rateLimit, burst)
}
