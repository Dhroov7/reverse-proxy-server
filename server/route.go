package server

import (
	configreader "reverse-proxy-server/configReader"
	proxyserver "reverse-proxy-server/proxyServer"

	"github.com/gin-gonic/gin"
)

func initRoutes() *gin.Engine {
	app := gin.Default()

	config := configreader.ReadConfig()
	newProxyServer := proxyserver.InitProxyServer(config)

	app.Any("/*any", newProxyServer.ServerRequest)

	return app
}
