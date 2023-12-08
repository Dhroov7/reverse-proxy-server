package configreader

import (
	"encoding/json"
	"fmt"
	"os"
)

type TargetServer struct {
	RoutePath  string `json:"routePath"`
	HealthPath string `json:"healthPath"`
}

type Config struct {
	RateLimiting struct {
		MaxRequestsPerMinute int `json:"maxRequestsPerMinute"`
		TokenRefillInterval  int `json:"tokenRefillInterval"`
	} `json:"rateLimiting"`
	TargetServers map[string]TargetServer `json:"targetServers"`
}

func ReadConfig() Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		os.Exit(1)
	}

	// Parse the JSON configuration
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error parsing configuration:", err)
		os.Exit(1)
	}
	fmt.Println(config)
	return config
}
