package rateLimiter

import (
	"errors"
	"net/http"
	configreader "reverse-proxy-server/configReader"
	"reverse-proxy-server/util"
	"strings"
	"time"
)

type TargetServersHealthConfig struct {
	Alive    bool
	Endpoint string
}

type RateLimiter struct {
	TargetServersHealth  map[string]TargetServersHealthConfig
	TargetServers        []string
	MaxRequestsPerMinute int
	TokenRefillInterval  int
	LastRefill           int
	Tokens               map[string]int
	CurrentServerIndex   int
	PathRouting          map[string]string
}

func InitRateLimiter(config configreader.Config) *RateLimiter {
	newRateLimiter := RateLimiter{
		TargetServersHealth:  make(map[string]TargetServersHealthConfig),
		MaxRequestsPerMinute: config.RateLimiting.MaxRequestsPerMinute,
		TokenRefillInterval:  config.RateLimiting.TokenRefillInterval,
		LastRefill:           util.GetCurrentTimeInMilli(),
		Tokens:               make(map[string]int),
		CurrentServerIndex:   1,
		PathRouting:          make(map[string]string),
	}

	for host, targetServerConfig := range config.TargetServers {
		newRateLimiter.Tokens[host] = newRateLimiter.MaxRequestsPerMinute
		newRateLimiter.TargetServersHealth[host] = TargetServersHealthConfig{
			Alive:    true,
			Endpoint: targetServerConfig.HealthPath,
		}
		newRateLimiter.TargetServers = append(newRateLimiter.TargetServers, host)
		newRateLimiter.PathRouting[targetServerConfig.RoutePath] = host
	}

	go newRateLimiter.HealthCheckPeriodically(2 * time.Second)

	return &newRateLimiter
}

func (rl *RateLimiter) GetServer(requestPath string) (string, error) {
	rl.refillTokens()
	selectedServer := rl.selectNextServer(requestPath)

	if rl.Tokens[selectedServer] > 0 {
		rl.Tokens[selectedServer]--
		return selectedServer, nil
	} else {
		return "", errors.New("Too many requests")
	}
}

func (rl *RateLimiter) refillTokens() {
	currentTime := util.GetCurrentTimeInMilli()

	timeElapsed := currentTime - rl.LastRefill
	tokenToRefill := (timeElapsed * (rl.TokenRefillInterval / 60 / 1000))

	for _, server := range rl.TargetServers {
		rl.Tokens[server] = util.Min(rl.MaxRequestsPerMinute, rl.Tokens[server]+tokenToRefill)
	}

	rl.LastRefill = currentTime
}

func (rl *RateLimiter) selectNextServer(requestPath string) string {
	var selectedServer string

	if len(rl.PathRouting) > 0 {
		for path, server := range rl.PathRouting {
			pathParts := strings.Split(path, "/")

			if strings.Contains(requestPath, strings.Join(pathParts[:len(pathParts)-1], "/")) && rl.TargetServersHealth[server].Alive {
				selectedServer = server
			}
		}
	} else {
		selectedServer = rl.TargetServers[rl.CurrentServerIndex]
		rl.CurrentServerIndex = (rl.CurrentServerIndex + 1) % len(rl.TargetServers)
	}

	return selectedServer
}

func (rl *RateLimiter) HealthCheckPeriodically(interval time.Duration) {
	for {
		for _, serverUrl := range rl.TargetServers {
			healthCheckRoute := serverUrl + rl.TargetServersHealth[serverUrl].Endpoint
			_, err := http.Get(healthCheckRoute)
			if err != nil {
				rl.updateServerHealth(serverUrl, false)
			} else {
				rl.updateServerHealth(serverUrl, true)
			}
		}
		time.Sleep(interval)
	}
}

func (rl *RateLimiter) updateServerHealth(serverHost string, healthStatus bool) {
	targetServerHealth := rl.TargetServersHealth[serverHost]
	targetServerHealth.Alive = healthStatus
	rl.TargetServersHealth[serverHost] = targetServerHealth
}
