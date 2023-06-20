package ctrl

import (
	"encoding/json"
	"fine-grained-openai-proxy/conf"
	"fine-grained-openai-proxy/svc"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"sort"
	"strconv"
)

type OpenAIModel struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
}

/*
InitModels re-fetch all available models from OpenAI, and compare them with db.
If there are new models, truncate db and insert new models.

Path:

	/admin/model/init

Args:

	GET auth: Admin Token
	POST id: OpenAI ApiKey ID
*/
func InitModels(c *fiber.Ctx) error {
	ids := c.FormValue("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Invalid id: " + err.Error()})
	}
	apikeySvc := svc.ApiKeySvc{}
	apikey, err := apikeySvc.ByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Invalid id: " + err.Error()})
	}
	// Create a new HTTP request
	req, err := http.NewRequest("GET", conf.OpenAIAPIBaseUrl+"/v1/models", nil)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Fetch models failed: " + err.Error()})
	}
	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+apikey.Key)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Fetch models failed: " + err.Error()})
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			return
		}
	}(resp.Body)

	// Process the OpenAI API response
	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		// Print the response body
		fmt.Println(string(body))
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Parse models failed: " + err.Error()})
	}

	oai := data.(map[string]interface{})
	if oai["error"] != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Fetch models failed: " + oai["error"].(string)})
	}

	// models = (oai["data"][i]["id"], oai["data"][i]["created"])
	var models []OpenAIModel
	for _, v := range oai["data"].([]interface{}) {
		model := v.(map[string]interface{})
		models = append(models, OpenAIModel{ID: model["id"].(string), Created: int64(model["created"].(float64))})
	}

	// sort models by created
	sort.Slice(models, func(i, j int) bool {
		return models[i].Created < models[j].Created
	})

	// compare models with db
	modelSvc := svc.ModelSvc{}
	dbModels, err := modelSvc.All()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Read db models failed: " + err.Error()})
	}

	// if there are new models (len changed or not match), truncate db and insert new models
	if len(models) != len(dbModels) || models[len(models)-1].ID != dbModels[len(dbModels)-1].Name {
		if err := modelSvc.Truncate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Truncate db models failed: " + err.Error()})
		}
		for _, v := range models {
			if err := modelSvc.Insert(v.ID); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(Resp{Code: 1, Error: "Insert db models failed: " + err.Error()})
			}
		}
		return c.Status(fiber.StatusOK).JSON(Resp{Code: 0, Data: "Init models success"})
	}
	return c.Status(fiber.StatusOK).JSON(Resp{Code: 0, Data: "Init models success, No new models"})
}
