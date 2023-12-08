package util

import (
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCurrentTimeInMilli() int {
	currentTime := time.Now()

	// Convert the current time to milliseconds
	currentTimeInMilli := currentTime.UnixNano() / int64(time.Millisecond)

	return int(currentTimeInMilli)
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SetupRequest(c *gin.Context, targetUrl *url.URL) {
	c.Request.Host = targetUrl.Host
	c.Request.URL.Host = targetUrl.Host
	c.Request.URL.Scheme = targetUrl.Scheme
	c.Request.URL.Path = c.Param("any")

	c.Request.Header.Set("X-Forwarded-For", c.ClientIP())
	c.Request.Header.Set("X-Forwarded-Host", c.Request.Header.Get("Host"))
	c.Request.Header.Set("Host", targetUrl.Host)
}
