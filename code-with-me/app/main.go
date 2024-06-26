package main

import (
	"code-with-me/internal/handlers"
	"code-with-me/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error loading env file : " + err.Error())
		return
	}
	codeService := service.New()
	codeHandlers := handlers.New(codeService)

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	router.Use(cors.New(config))
	router = handlers.InitRouter(router, codeHandlers)
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "9989"
		slog.Info("no port found : starting application of default : 9989")
	}
	router.Run(":" + port)
	//router.Run()
}
