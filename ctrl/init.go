package ctrl

import (
	"fine-grained-openai-proxy/conf"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code  int         `json:"code"` // 0: success, 1: failed
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg,omitempty"`   // success message
	Error string      `json:"error,omitempty"` // failed message
}

// AuthMiddleware is a middleware to check _GET["auth"] == conf.AdminToken
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.Query("auth")
        if auth != conf.AdminToken {
            c.JSON(http.StatusUnauthorized, Resp{Code: 1, Error: "auth failed"})
            c.Abort()
            return
        }
        c.Next()
    }
}

func InitRouter(router *gin.Engine, allowAdmin ...bool) {
    router.Any("/v1/*any", gin.WrapF(Proxy))

    if !allowAdmin[0] {
        return
    }
    admin := router.Group("/admin")
    admin.Use(AuthMiddleware())
    // ApiKey
    admin.GET("/apikey/all", AllApiKeys)
    admin.POST("/apikey/insert", InsertApiKey)
    admin.POST("/apikey/delete", DeleteApiKey)
    // FineGrainedKey
    admin.GET("/fgkey/all", AllFineGrainedKeys)
    admin.POST("/fgkey/parentid", FineGrainedKeysByParentID)
    admin.POST("/fgkey/insert", InsertFineGrainedKey)
    admin.POST("/fgkey/Update", UpdateFineGrainedKey)
    admin.POST("/fgkey/delete", DeleteFineGrainedKey)
    // Model
    admin.GET("/model/all", AllModels)
    admin.POST("/model/init", InitModels)
}