package videos

import (
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterRoutes(basePath string, r *gin.Engine, m *mongo.Client) {
	s := newVideoService(m)
	c := newVideoController(s)

	r.GET(basePath+"/videos", c.ListVideos)
	r.POST(basePath+"/videos", c.CreateVideo)
	r.DELETE(basePath+"/videos/:id", c.DeleteVideo)
	r.GET(basePath+"/videos/:id", c.GetVideo)
	r.PATCH(basePath+"/videos/:id", c.UpdateVideo)
	r.PUT(basePath+"/videos/:id", c.UpdateVideo)
}
