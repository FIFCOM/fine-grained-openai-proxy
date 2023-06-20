package ctrl

import (
	"encoding/json"
	"fine-grained-openai-proxy/conf"
	"fine-grained-openai-proxy/svc"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"strings"
	"time"
)

/*
OpenAIAPIHandler proxy all requests to OpenAI API

Proxy Path: /v1/*

 1. Get request Header "Bearer <fine-grained_key>".
 2. If fine-grained_key is valid, get POST body.
 3. If _POST["model"] is in the list of allowed models,
 4. change Header "Bearer <fine-grained_key>" to "Bearer <real_openai_key>"
 5. proxy request to OpenAI API.
*/
func OpenAIAPIHandler(c *fiber.Ctx) error {
	// 1. Get request Header "Bearer <fine-grained_key>".
	authHeader := c.Get("Authorization")
	var token string
	fgSvc := svc.FineGrainedKeySvc{}
	var fgKey *svc.FineGrainedKey
	fgKeyOK := false
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		} else {
			Proxy(c) // no token
			return nil
		}
		var err error
		fgKey, err = fgSvc.ByHash(svc.Hash(token))
		timeNow := time.Now().Unix()
		if err != nil && fgKey.RemainCalls > 0 && fgKey.Expire <= timeNow {
			Proxy(c) // invalid token
			return err
		}
		fgKeyOK = true
	} else {
		Proxy(c) // No Authorization header
		return nil
	}

	var reqData interface{}
	if err := c.BodyParser(&reqData); err != nil {
		Proxy(c)
		return err
	}

	params := reqData.(map[string]interface{})
	models, ok := params["model"]
	if !ok {
		Proxy(c) // no model
		return nil
	}
	// println(models.(string))

	modelSvc := svc.ModelSvc{}
	model, err := modelSvc.ByName(models.(string))
	if err != nil {
		Proxy(c) // model not in the list of allowed models
		return err
	}
	modelOK := false
	var arr []int64
	// json_decode fgKey.List to arr
	_ = json.Unmarshal([]byte(fgKey.List), &arr)
	if fgKey.Type == "whitelist" {
		for _, mid := range arr {
			if mid == model.ID {
				modelOK = true
			}
		}
	} else if fgKey.Type == "blacklist" {
		modelOK = true
		for _, mid := range arr {
			if mid == model.ID {
				modelOK = false
			}
		}
	}

	if !fgKeyOK || !modelOK {
		Proxy(c)
		return nil
	}
	fgKey.RemainCalls -= 1
	_ = fgSvc.Update(fgKey)
	// c.Response().Header.Del(fiber.HeaderServer)

	// change Authorization header
	apikeySvc := svc.ApiKeySvc{}
	apiKey, _ := apikeySvc.ByID(fgKey.ParentID)
	// unset Authorization header
	c.Request().Header.Del("Authorization")
	c.Request().Header.Set("Authorization", "Bearer "+apiKey.Key)
	Proxy(c)
	return nil
}

// Proxy ctx request to OpenAI API
func Proxy(c *fiber.Ctx) {
	proxyUrl := fmt.Sprintf("%s/v1/%s", conf.OpenAIAPIBaseUrl, c.Params("*", ""))

	if err := proxy.Do(c, proxyUrl); err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(Resp{
			Code:  1,
			Error: "Proxy failed. + " + err.Error(),
		})
	}
	c.Response().Header.Del(fiber.HeaderServer)
	// print ctx Header and Body
	println(c.Request().Header.String())
	println(string(c.Request().Body()))
}
