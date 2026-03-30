package videos

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type videoController struct {
	service videoService
}

func newVideoController(s videoService) *videoController {
	return &videoController{service: s}
}

func (vc *videoController) GetVideo(c *gin.Context) {
	videoID := c.Param("id")
	video, err := vc.service.GetVideoByID(videoID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	c.JSON(http.StatusOK, video)
}

func (vc *videoController) ListVideos(c *gin.Context) {
	videos, err := vc.service.ListVideos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list videos"})
		return
	}
	c.JSON(http.StatusOK, videos)
}

func (vc *videoController) CreateVideo(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	videoID, err := vc.service.CreateVideo(req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": videoID})
}

func (vc *videoController) DeleteVideo(c *gin.Context) {
	videoID := c.Param("id")

	if err := vc.service.DeleteVideo(videoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete video"})
		return
	}

	c.Status(http.StatusOK)
}

func (vc *videoController) UpdateVideo(c *gin.Context) {
	videoID := c.Param("id")
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	isUpsert := true
	if c.Request.Method == http.MethodPut {
		isUpsert = true
	}

	if err := vc.service.UpdateVideo(videoID, req.Path, isUpsert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update video"})
		return
	}

	c.Status(http.StatusOK)
}
