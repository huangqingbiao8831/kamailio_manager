package api

import (
	"kamailio-manager/internal/client"
	"kamailio-manager/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MtreeReloadRequestParam struct {
	Tname string `json:"mtreeName"` //对应
}

func HandleMtreeReload(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MtreeReloadRequestParam
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"msg": "无效的参数格式", "code": 400})
			logger.Log.Warn("Invalid dialplan parameters", zap.Error(err))
			return
		}
		method := "mtree.reload"
		var params []interface{}
		if req.Tname != "" {
			params = append(params, req.Tname)
			logger.Log.Info("params:", zap.String("mtreeName", req.Tname))
		}
		result, err := kClient.Invoke(method, params...)
		if err != nil {
			c.JSON(500, gin.H{"msg": err.Error(), "code": 500})
			logger.Log.Warn("Invoke rpc error", zap.Error(err))
			return
		}

		c.JSON(200, gin.H{
			"code":           200,
			"msg":            "操作成功",
			"result":         result,
			"reloaded_mtree": req.Tname,
		})
	}
}
