package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kamailio-manager/internal/client"
	"kamailio-manager/internal/logger"
	//"net/http"
)

func HandleDispatcherReload(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := kClient.Invoke("dispatcher.reload")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "code": 500})
			logger.Log.Warn("Enter HandleDispatcherReload,get error:", zap.Error(err))
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "操作成功", "data": res})
	}
}
