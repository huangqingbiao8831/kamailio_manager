package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kamailio-manager/internal/client"
	"kamailio-manager/internal/logger"
	//"net/http"
)

func HandleDlgStatsActive(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := kClient.Invoke("dlg.stats_active")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			logger.Log.Warn("Enter HandleDlgStatsActive,get error:",zap.Error(err))
			return
		}
		c.JSON(200, gin.H{"status": "success", "data": res})
	}
}
