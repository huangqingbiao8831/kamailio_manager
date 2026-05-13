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

type MtreeCommonRequestParam struct {
	Method   string `json:"mtreeMethod" binding:"required"`
	TreeRoot string `json:"treeRootName" binding:"required"` //对应table名字
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

func HandleMtreeCommonCommand(kClient *client.KamailioJSONClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 参数绑定与校验
		var req MtreeCommonRequestParam
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Log.Warn("Mtree 参数绑定失败", zap.Error(err))
			c.JSON(400, gin.H{"code": 400, "error": "无效的参数格式，treeRootName 和 mtreeMethod 必填"})
			return
		}

		// 2. 转换并校验 Method
		var rpcMethod string
		switch req.Method {
		case "reload", "summary", "list":
			rpcMethod = "mtree." + req.Method
		default:
			// 如果是自定义匹配逻辑，可以扩展 mtree.match
			c.JSON(400, gin.H{"code": 400, "error": "不支持的 mtree 操作方法"})
			return
		}

		// 3. 树根名称映射 (白名单模式)
		// 对应你在 Kamailio 配置文件中定义的 dbtable 名称或 tname
		actualTreeName := ""
		if req.TreeRoot == "" {
			// 允许直接传入已知的合法名称
			logger.Log.Warn("收到空的 mtree 树根名称")
			c.JSON(400, gin.H{"msg": "treeRootName 必须提供", "code": 400})
			return
		} else {
			actualTreeName = req.TreeRoot
		}

		// 4. 构建参数列表
		// mtree.reload 和 mtree.list 通常只需要 tname 这一个参数
		params := []interface{}{actualTreeName}

		// 5. 执行 Kamailio RPC 调用
		logger.Log.Info("执行 mtree RPC",
			zap.String("method", rpcMethod),
			zap.String("tree", actualTreeName))

		result, err := kClient.Invoke(rpcMethod, params...)
		if err != nil {
			logger.Log.Error("Kamailio mtree RPC 调用失败", zap.Error(err), zap.String("method", rpcMethod))
			c.JSON(500, gin.H{"code": 500, "error": "后端通讯失败: " + err.Error()})
			return
		}

		// 6. 成功响应
		c.JSON(200, gin.H{
			"code":      200,
			"msg":       "操作成功",
			"data":      result,
			"tree-name": actualTreeName,
		})
	}
}
