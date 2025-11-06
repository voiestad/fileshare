package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", index)
	router.GET("/download/:id", getFile)
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.Run("0.0.0.0:8080")
}

func index(c *gin.Context) {
	availableFiles, _ := getAvailableFiles()
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":          "fileshare",
		"availableFiles": availableFiles,
	})
}

type File struct {
	Name string
	Id   int
}

func getAvailableFiles() ([]File, error) {
	availableFiles, err := os.ReadDir("files")
	if err != nil {
		return make([]File, 0), err
	}
	files := make([]File, len(availableFiles))
	for i, file := range availableFiles {
		files[i] = File{Name: file.Name(), Id: i}
	}
	return files, nil
}

func getFile(c *gin.Context) {
	inputFile := c.Param("id")
	if inputFile == "" {
		c.String(http.StatusBadRequest, "Missing 'file' parameter")
		return
	}
	fileIdx, err := strconv.Atoi(inputFile)
	if err != nil {
		c.String(http.StatusBadRequest, "'file' parameter should be an integer")
		return
	}
	availableFiles, err := os.ReadDir("files")
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not load the files")
		return
	}
	if fileIdx < 0 || fileIdx >= len(availableFiles) {
		c.String(http.StatusBadRequest, "'file' parameter not within legal range")
		return
	}
	fileName := availableFiles[fileIdx].Name()

	file := filepath.Join("files", fileName)
	c.File(file)
}
