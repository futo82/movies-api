package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"./controllers"
	"./db"
	"./middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitializeDynamoDB()

	router := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(middlewares.AuthMiddleware())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1/api")
	{
		v1.POST("/movies", controllers.CreateMovie)
		v1.GET("/movies/:id", controllers.RetrieveMovie)
		v1.PUT("/movies/:id", controllers.UpdateMovie)
		v1.DELETE("/movies/:id", controllers.DeleteMovie)
	}

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to startup server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server is shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Encountered an error while shutting down server: ", err)
	}
	log.Println("Server shutdown")
}
