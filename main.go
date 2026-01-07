package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/voiestad/fileshare/internal/database"
)

var fileDirName string = "files"

func main() {
	os.Mkdir(fileDirName, 0775)
	err := database.Init()
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", index)
	router.GET("/download/:id", getFile)
	router.POST("/upload", addFiles)
	router.POST("/delete/:id", deleteFile)
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.StaticFile("/style.css", "./static/style.css")
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
	Id   uuid.UUID
}

func getAvailableFiles() ([]File, error) {
	availableFiles, err := database.GetFiles()
	if err != nil {
		return make([]File, 0), err
	}
	files := make([]File, 0)
	for availableFiles.Next() {
		var id, name, time string
		if err := availableFiles.Scan(&id, &name, &time); err != nil {
			log.Println("scan error:", err)
			continue
		}
		parsedId, err := uuid.Parse(id)
		if err != nil {
			log.Println("scan error:", err)
			continue
		}
		files = append(files, File{Name: name, Id: parsedId})
	}
	availableFiles.Close()
	return files, nil
}

func getFile(c *gin.Context) {
	inputFile := c.Param("id")
	if inputFile == "" {
		c.String(http.StatusBadRequest, "Missing 'file' parameter")
		return
	}
	id, err := uuid.Parse(inputFile)
	if err != nil {
		c.String(http.StatusBadRequest, "'file' parameter should be a valid UUID")
		return
	}
	res := database.GetFile(id)
	var fileName string
	err = res.Scan(&fileName)
	if err != nil {
		c.String(http.StatusNotFound, "Could not find file")
		return
	}
	file := filepath.Join(fileDirName, id.String())
	c.File(file)
}

type AddFilesResponse struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
}

func addFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	files := form.File["files"]
	success := make([]string, 0)
	failed := make([]string, 0)
	if err != nil {
		c.String(http.StatusBadRequest, "Could not get file: %v", err)
		return
	}
	for _, file := range files {
		id, err := database.AddFile(file)
		if err != nil {
			failed = append(failed, file.Filename)
			continue
		}
		err = c.SaveUploadedFile(file, filepath.Join(fileDirName, id.String()))
		if err != nil {
			database.RemoveFile(id)
			failed = append(failed, file.Filename)
			continue
		}
		success = append(success, file.Filename)
	}
	c.JSON(http.StatusOK, AddFilesResponse{Success: success, Failed: failed})
}

func deleteFile(c *gin.Context) {
	inputFile := c.Param("id")
	if inputFile == "" {
		c.String(http.StatusBadRequest, "Missing 'file' parameter")
		return
	}
	id, err := uuid.Parse(inputFile)
	if err != nil {
		c.String(http.StatusBadRequest, "'file' parameter should be a valid UUID")
		return
	}
	file := filepath.Join(fileDirName, id.String())
	err = os.Remove(file)
	if err != nil {
		c.String(http.StatusBadRequest, "File could not be deleted")
		return
	}
	database.RemoveFile(id)
	c.String(http.StatusOK, "File deleted successfully!")
}
