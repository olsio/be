package main

import (
	"archive/zip"
	"encoding/json"
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
)

func getMeasurements(c *gin.Context) {
	files, _ := ioutil.ReadDir(dataDirectoryPath)

	var measurements []Measurement
	for _, file := range files {
		jsonFile, _ := ioutil.ReadFile(path.Join(dataDirectoryPath, file.Name()))
		var jsonObject DeviceContent
		json.Unmarshal(jsonFile, &jsonObject)
		measurements = append(measurements, jsonObject.Measurements...)
	}
	c.JSON(http.StatusOK, measurements)
}

func buildZip(c *gin.Context) {
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	zipit(dataDirectoryPath, tmpfile.Name())
	c.File(tmpfile.Name())
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
