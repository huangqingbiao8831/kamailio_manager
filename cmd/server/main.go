package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kamailio-manager/internal/api"
	"kamailio-manager/internal/client"
	"kamailio-manager/internal/config"
	"kamailio-manager/internal/logger"
	"time"
)

// PollKamailioStatus 持续探测 Kamailio 状态直到成功
func PollKamailioStatus(kClient *client.KamailioJSONClient) {
	// 给 Kamailio 留出基础启动缓冲时间
	time.Sleep(2 * time.Second)

	logger.Log.Info("开始启动后自动探测 Kamailio 业务就绪状态...")

	// 设置 3 秒一次的定时探测
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 调用之前定义的原生 HTTP 请求方法
			res, err := kClient.InvokeNativeHTTP("/api/startSuccessful", map[string]string{
				"action": "init_report",
			})

			if err == nil {
				logger.Log.Info("Kamailio 业务接口调用成功!", zap.String("response", res))
				return // 成功后退出循环，结束协程
			}

			logger.Log.Warn("Kamailio 尚未就绪或接口报错，准备重试", zap.Error(err))
		}
	}
}

func main() {
	// 1. Initialize Log
	logger.Init()
	defer logger.Log.Sync()

	// 2. Load Config
	conf := config.LoadConfig()

	// 3. Initialize Clients
	kClient := &client.KamailioJSONClient{
		Endpoint: conf.Kamailio.RPCUrl,
		Timeout:  time.Duration(conf.Kamailio.Timeout) * time.Second,
	}
	sClient := &client.SupervisorClient{
		Endpoint: conf.Supervisor.Url,
	}

	// 4. Setup Router
	r := gin.New()
	// ... (Rest of your middleware and routes)
	r.Use(gin.Recovery())

	// 日志中间件
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Log.Info("HTTP请求",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("latency", time.Since(start)),
		)
	})
	//可以增加一些全局的中间件
	v1 := r.Group("/v1")
	//v1.Use(fun_mid_A(),fun_mid_B(),fun_mid_C())
	{
		service := v1.Group("/service")
		//service.Use(AuthMiddleware(), RateLimitMiddleware(), TimeWindowMiddleware())
		{
			service.POST("/kamailio/reload", func(c *gin.Context) {
				// 执行之前定义的逻辑，如重载限制表
				res, err := kClient.Invoke("rpc.reload_restrictions")
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"status": "success", "data": res})
			})

			service.POST("/dispatcher/reload", api.HandleDispatcherReload(kClient))

			service.POST("/dialplan/reload", api.HandleDialplanReload(kClient))

			service.POST("/dialplan/translate", api.HandleDialplanTranslate(kClient))

			service.POST("/permission/address-reload", api.HandlePermissionAdressReload(kClient))

			service.POST("/tm/stats", api.HandleTmModStats(kClient))

			service.POST("/dialog/stats_active", api.HandleDlgStatsActive(kClient))
		}
		// 2. 系统级控制 (Supervisor XML-RPC)
		system := v1.Group("/system")
		//system.Use(AuthMiddleware(), RateLimitMiddleware(), TimeWindowMiddleware())
		{
			system.POST("/process/stop", func(c *gin.Context) {
				pName := conf.Supervisor.ProcessName
				logger.Log.Warn("start stopping ", zap.String("process", pName))

				//stopping
				success, err := sClient.ControlProcess("stop", pName)

				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"status": "process stop", "success": success})
			})

			// 3. process start
			system.POST("/process/start", func(c *gin.Context) {
				pName := conf.Supervisor.ProcessName
				logger.Log.Warn("start starting ", zap.String("process", pName))

				//stopping
				success, err := sClient.ControlProcess("start", pName)

				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"status": "process start", "success": success})
			})
		}
	}

	logger.Log.Info("管理服务启动", zap.String("port", conf.Server.Port))
	//增加对kamailio启动状态检测，并触发kamailio通过rabbitmq发送消息到后端
	go PollKamailioStatus(kClient)
	r.Run(":" + conf.Server.Port)
}
