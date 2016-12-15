package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsjamie/gin-cors"
)

const (
	defaultPort       = "5000"
	dataDirectoryPath = "./data"
)

type measurement struct {
	ID       string `json:"uuid"`
	TrialID  string `json:"trialId"`
	Subject  string `json:"subject"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Duration string `json:"duration"`
	Target   string `json:"target"`
	Response string `json:"response"`
	Correct  bool   `json:"correct"`
}

func main() {
	port := getPort()

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
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/measurements", func(c *gin.Context) {
		tmpfile, err := ioutil.TempFile("", "")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		zipit(dataDirectoryPath, tmpfile.Name())
		c.File(tmpfile.Name())
	})

	initializeDataDirectory(dataDirectoryPath)

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

func saveMeasurements(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": err,
		})
		return
	}

	fmt.Println(string(body))
	ioutil.WriteFile(generateFileName(), body, 0600)
	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
	})
}

func initializeDataDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return err
}

func generateFileName() string {
	uniqueID, _ := uuid.NewRandom()
	filename := fmt.Sprintf("%s %s.txt", time.Now().Format("2006-01-02 150405"), uniqueID)
	target := path.Join(dataDirectoryPath, filename)
	return target
}

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
