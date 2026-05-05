package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *ApiServer) UploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	links := make([]string, 0, len(files))
	for _, file := range files {
		link, err := s.fileClient.UploadPhoto(c, file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		links = append(links, link)
	}
	c.JSON(http.StatusOK, gin.H{"links": links})
}
