package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

const (
	defaultPort       = "5000"
	dataDirectoryPath = "./data"
)

func main() {
	port := getPort()

	initializeDataDirectory(dataDirectoryPath)

	router := gin.New()
	corsConfig := cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}

	router.Use(cors.Middleware(corsConfig))
	router.Use(gin.Logger())
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.StaticFile("/asset-manifest.json", "./static/asset-manifest.json")
	router.Static("/static", "./static/static")

	router.GET("/zip", buildZip)
	router.GET("/measurements", getMeasurements)
	router.POST("/measurements", saveMeasurements)
	router.Run(":" + port)
}

func getPort() string {
	port := os.Getenv("PORT")

	if port != "" {
		return port
	}

	if len(os.Args) != 2 {
		log.Print("Using default port " + defaultPort)
		return defaultPort
	}

	return os.Args[1]
}

func initializeDataDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return err
}
