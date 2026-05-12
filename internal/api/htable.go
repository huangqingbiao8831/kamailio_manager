package api

import (
	"kamailio-manager/internal/client"
	"kamailio-manager/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HtableReloadRequestParam struct {
	Table string `json:"htableName"` //对应
}

// 重新加载htable表中
func HandleHtableReload(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req HtableReloadRequestParam

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"msg": "无效的参数格式", "code": 400})
			logger.Log.Warn("Invalid dialplan parameters", zap.Error(err))
			return
		}
		method := "htable.reload"
		var params []interface{}
		if req.Table != "" {
			logger.Log.Info("get htable tableName", zap.String("tableName", req.Table))
			switch req.Table {
			case "k_htable_num_pools":
				params = append(params, "barringTimeMap")
			case "k_htable_time_rules":
				params = append(params, "numPoolMap")
			default:
				logger.Log.Error("Unknown tableName", zap.String("tableName", req.Table))
			}
		}
		result, err := kClient.Invoke(method, params...)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			logger.Log.Warn("Invoke rpc error", zap.Error(err))
			return
		}

		c.JSON(200, gin.H{
			"code":            200,
			"msg":             "操作成功",
			"result":          result,
			"reloaded_htable": req.Table,
		})
	}
}
