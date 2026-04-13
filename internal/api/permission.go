package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kamailio-manager/internal/client"
	"kamailio-manager/internal/logger"
	//"net/http"
)

func HandlePermissionAdressReload(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := kClient.Invoke("permissions.addressReload")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			logger.Log.Warn("Enter HandlePermissionAdressReload,get error:",zap.Error(err))
			return
		}
		c.JSON(200, gin.H{"status": "success", "data": res})
	}
}
