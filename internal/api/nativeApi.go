package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kamailio-manager/internal/client"
	"kamailio-manager/internal/logger"
	"net/http"
)

// HandleKamailioNativeAPI 处理通用的原生 HTTP 转发请求
// 对应调用: POST /v1/service/kamailio/api?path=/api/xxx
func HandleKamailioNativeAPI(c *gin.Context, kClient *client.KamailioJSONClient) {
	// 1. 获取路径参数
	apiPath := c.Query("path")
	if apiPath == "" {
		logger.Log.Warn("Native API 调用缺失 path 参数")
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'path' is required"})
		return
	}

	// 2. 提取并透传所有 Query 参数
	queryParams := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if k != "path" && len(v) > 0 {
			queryParams[k] = v[0]
		}
	}

	// 3. 调用 Kamailio 的 xhttp 接口 InvokeNativeHTTP
	res, err := kClient.InvokeNativeHTTP(apiPath, queryParams)
	if err != nil {
		logger.Log.Error("Invoke Native HTTP Failed", zap.String("path", apiPath), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"path":   apiPath,
		"data":   res,
	})
}

// GetNativeAPIHandler 返回标准的 gin.HandlerFunc 闭包
func GetNativeAPIHandler(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleKamailioNativeAPI(c, kClient)
	}
}
