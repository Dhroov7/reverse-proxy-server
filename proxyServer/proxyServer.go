package proxyserver

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	configreader "reverse-proxy-server/configReader"
	"reverse-proxy-server/logger"
	"reverse-proxy-server/rateLimiter"
	"reverse-proxy-server/util"
	"time"

	"github.com/gin-gonic/gin"
)

type TargetServersHealthConfig struct {
	Alive    bool
	Endpoint string
}

type ProxyServer struct {
	rateLimiter rateLimiter.RateLimiter
}

func InitProxyServer(config configreader.Config) *ProxyServer {
	newProxyServer := ProxyServer{
		rateLimiter: *rateLimiter.InitRateLimiter(config),
	}

	return &newProxyServer
}

func (ps *ProxyServer) ServerRequest(c *gin.Context) {
	requestPath := c.Param("any")

	selectedServer, err := ps.rateLimiter.GetServer(requestPath)

	if len(selectedServer) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Resource Not Found!"})
		return
	}

	if err != nil {
		c.IndentedJSON(http.StatusTooManyRequests, gin.H{"message": "Too many Requests"})
		return
	}

	var targetURL, _ = url.Parse(selectedServer)
	var proxy = httputil.NewSingleHostReverseProxy(targetURL)

	util.SetupRequest(c, targetURL)
	proxy.ServeHTTP(c.Writer, c.Request)

	logger.Log(fmt.Sprintf("%s %s", time.Now().UTC(), targetURL))
}
