package su2

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const uploadsPath = "su2/uploads"

func GinFileServer() {
	flag.StringVar(&port, "p", "8080", "port to listen on")
	flag.StringVar(&address, "a", "localhost", "address to listen on")

	flag.Parse()

	router := gin.Default()

	server := &http.Server{
		Addr:              net.JoinHostPort(address, port),
		Handler:           router,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 2,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	RunGinFileServer(router, server)
}

func RunGinFileServer(router *gin.Engine, s *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	router.MaxMultipartMemory = 8 << 20
	router.POST("/files", HandleUploadFile)
	router.GET("/files/:id", HandleDownloadFile)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			return
		}
	}()

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s\n", err.Error())
	}

	log.Printf("Server exiting\n")
}

func HandleUploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file.Filename = uuid.NewString() + filepath.Ext(file.Filename)

	err = c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", uploadsPath, file.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error saving file": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file ID": file.Filename})
}

func HandleDownloadFile(c *gin.Context) {
	id := c.Param("id")

	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	filePath := filepath.Join(uploadsPath, id)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.File(filePath)
}
