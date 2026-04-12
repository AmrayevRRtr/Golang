package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	limitDuration = 10 * time.Second
	maxRequests   = 5
)

type clientRequest struct {
	count      int
	lastAccess time.Time
}

var (
	clients = make(map[string]*clientRequest)
	mu      sync.Mutex
)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		userID, exists := c.Get("userID")
		if exists {
			key = userID.(string)
		} else {
			key = c.ClientIP()
		}

		mu.Lock()

		now := time.Now()
		client, found := clients[key]

		if !found || now.Sub(client.lastAccess) > limitDuration {
			clients[key] = &clientRequest{
				count:      1,
				lastAccess: now,
			}
			mu.Unlock()
			c.Next()
			return
		}

		if client.count >= maxRequests {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			return
		}

		client.count++
		client.lastAccess = now
		mu.Unlock()

		c.Next()
	}
}
