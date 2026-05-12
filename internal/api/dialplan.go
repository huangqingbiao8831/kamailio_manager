package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kamailio-manager/internal/client"
	"kamailio-manager/internal/logger"
	"net/http"
)

type DialplanTranslateRequest struct {
	DPID   int    `json:"dpid" binding:"required"`
	SrcStr string `json:"srcstr" binding:"required"`
}

func HandleDialplanTranslate(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DialplanTranslateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Log.Warn("Invalid dialplan parameters", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing fields dpid(int) and srcstr(string)"})
			return
		}

		res, err := kClient.Invoke("dialplan.translate", req.DPID, req.SrcStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": res})
	}
}

// dialplan reload api
func HandleDialplanReload(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := kClient.Invoke("dispatcher.reload")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "success", "data": res})
	}
}
