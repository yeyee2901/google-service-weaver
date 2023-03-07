package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yeyee2901/service-weaver/service"
)

func main() {
	s, err := service.NewService(gin.ReleaseMode)
	if err != nil {
		log.Fatalln(err)
	}

	s.RegisterRouting()
	errChan := s.Run("localhost:12345")
	for err := range errChan {
		log.Fatalln(err)
	}
}
