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

type HtableCommonRequestParam struct {
	Method   string `json:"htableMethod" binding:"required"`
	Table    string `json:"htableName" binding:"required"` //对应table名字
	HtabeKey string `json:"htableKey"`
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

func HandleHtableCommonComand(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req HtableCommonRequestParam
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"code": 400, "error": "无效的参数格式"})
			return
		}

		// 1. 转换并校验 Method
		// 确保前端传的是 "get", "reload", "delete" 等，并映射为 Kamailio RPC 指令
		var rpcMethod string
		switch req.Method {
		case "reload", "get","stats","dump","listTables":
			rpcMethod = "htable." + req.Method
		case "":
			rpcMethod = "htable.get" // 默认行为
		default:
			c.JSON(400, gin.H{"code": 400, "error": "不支持的操作方法"})
			return
		}

		// 2. 表名映射 (保持之前的白名单逻辑)
		actualTableName := ""
		switch req.Table {
		case "k_htable_num_pools":
			actualTableName = "numPoolMap"
		case "k_htable_time_rules":
			actualTableName = "barringTimeMap"
		default:
			c.JSON(400, gin.H{"code": 400, "error": "不支持的表名"})
			return
		}

		// 3. 动态构建参数
		params := []interface{}{actualTableName}

		// 核心逻辑：只有非 reload 操作且 Key 不为空时才追加 Key 参数
		if req.Method == "get" && req.HtabeKey != "" {
			params = append(params, req.HtabeKey)
		}

		// 4. 执行调用
		result, err := kClient.Invoke(rpcMethod, params...)
		if err != nil {
			logger.Log.Error("Kamailio RPC 失败", zap.String("method", rpcMethod), zap.Error(err))
			c.JSON(500, gin.H{"code": 500, "error": "执行失败: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"code":   200,
			"msg":    "操作成功",
			"data": result,
		})
	}
}
