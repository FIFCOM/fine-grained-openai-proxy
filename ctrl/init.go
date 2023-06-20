package ctrl

import (
	"fine-grained-openai-proxy/conf"
	"github.com/gofiber/fiber/v2"
)

type Resp struct {
	Code  int         `json:"code"` // 0: success, 1: failed
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg,omitempty"`   // success message
	Error string      `json:"error,omitempty"` // failed message
}

// AuthMiddleware is a middleware to check _GET["auth"] == conf.AdminToken
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Query("auth")
		if auth != conf.AdminToken {
			return c.Status(fiber.StatusUnauthorized).
				JSON(Resp{Code: 1, Error: "auth failed"})
		}
		return c.Next()
	}
}

func InitRouter(app *fiber.App, allowAdmin ...bool) {
	app.All("/v1/*", OpenAIAPIHandler)

	if !allowAdmin[0] {
		return
	}
	admin := app.Group("/admin")
	admin.Use(AuthMiddleware())
	// ApiKey
	admin.Get("/apikey/all", AllApiKeys)
	admin.Post("/apikey/insert", InsertApiKey)
	admin.Post("/apikey/delete", DeleteApiKey)
	// FineGrainedKey
	admin.Get("/fgkey/all", AllFineGrainedKeys)
	admin.Post("/fgkey/parentid", FineGrainedKeysByParentID)
	admin.Post("/fgkey/insert", InsertFineGrainedKey)
	admin.Post("/fgkey/Update", UpdateFineGrainedKey)
	admin.Post("/fgkey/delete", DeleteFineGrainedKey)
	// Model
	admin.Post("/model/init", InitModels)
}
