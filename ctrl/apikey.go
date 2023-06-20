package ctrl

import (
	"fine-grained-openai-proxy/svc"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

/*
AllApiKeys Get all OpenAI API keys

Path:

	/admin/apikey/all

Args:

	GET auth: Admin token
*/
func AllApiKeys(c *fiber.Ctx) error {
	apikeySvc := svc.ApiKeySvc{}
	apikeys, err := apikeySvc.All()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(Resp{Code: 1, Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(Resp{
		Code: 0,
		Data: apikeys,
	})
}

/*
InsertApiKey Insert a new OpenAI API key

Path:

	/admin/apikey/insert

Args:

	GET auth: Admin token
	POST key: OpenAI API key
*/
func InsertApiKey(c *fiber.Ctx) error {
	key := c.FormValue("key")
	apikeySvc := svc.ApiKeySvc{}
	err := apikeySvc.Insert(key)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(Resp{Code: 1, Error: "Insert API Key Error: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(Resp{
		Code: 0,
		Msg:  "Insert API key successfully",
	})
}

/*
DeleteApiKey Delete an OpenAI API key

Path:

	/admin/apikey/delete

Args:

	GET auth: Admin token
	POST id: OpenAI API key id
*/
func DeleteApiKey(c *fiber.Ctx) error {
	ids := c.FormValue("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(Resp{Code: 1, Error: "Delete API Key Error: " + err.Error()})
	}

	apikeySvc := svc.ApiKeySvc{}
	err = apikeySvc.Delete(&svc.ApiKey{ID: id})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(Resp{Code: 1, Error: "Delete API Key Error: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(Resp{
		Code: 0, Msg: "Delete API key successfully",
	})
}
