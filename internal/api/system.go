package api

import (
        "net/http"
        "kamailio-manager/internal/client"
        "kamailio-manager/internal/logger"
        "github.com/gin-gonic/gin"
        "go.uber.org/zap"
)

func HandleProcessControl(sClient *client.SupervisorClient, action string, pName string) gin.HandlerFunc {
        return func(c *gin.Context) {
                logger.Log.Warn("Attempting process control", zap.String("action", action), zap.String("process", pName))
                success, err := sClient.ControlProcess(action, pName)
                if err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                        return
                }
                c.JSON(http.StatusOK, gin.H{"status": "process " + action, "success": success})
        }
}

