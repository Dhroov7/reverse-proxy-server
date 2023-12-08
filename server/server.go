package server

import "fmt"

func InitServer() {
	gin := initRoutes()
	go func() {
		err := gin.Run("localhost:8080")
		if err != nil {
			fmt.Println("Error while starting server")
		}
	}()

	ch := make(chan string)
	<-ch
}
