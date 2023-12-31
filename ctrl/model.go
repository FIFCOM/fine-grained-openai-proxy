package ctrl

import (
	"encoding/json"
	"fine-grained-openai-proxy/conf"
	"fine-grained-openai-proxy/svc"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
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
func InitModels(c *gin.Context) {
    ids := c.PostForm("id")
    id, err := strconv.ParseInt(ids, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Invalid id: " + err.Error()})
        return
    }
    apikeySvc := svc.ApiKeySvc{}
    apikey, err := apikeySvc.ByID(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Invalid id: " + err.Error()})
        return
    }
    // Create a new HTTP request
    req, err := http.NewRequest("GET", conf.OpenAIAPIBaseUrl+"/v1/models", nil)
    if err != nil {
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Fetch models failed: " + err.Error()})
        return
    }
    // Set the Authorization header
    req.Header.Set("Authorization", "Bearer "+apikey.Key)

    // Send the HTTP request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Fetch models failed: " + err.Error()})
        return
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
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Parse models failed: " + err.Error()})
        return
    }

    oai := data.(map[string]interface{})
    if oai["error"] != nil {
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Fetch models failed: " + oai["error"].(string)})
        return
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
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Read db models failed: " + err.Error()})
        return
    }

    // if there are new models (len changed or not match), truncate db and insert new models
    if len(models) != len(dbModels) || models[len(models)-1].ID != dbModels[len(dbModels)-1].Name {
        if err := modelSvc.Truncate(); err != nil {
            c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Truncate db models failed: " + err.Error()})
            return
        }
        for _, v := range models {
            if err := modelSvc.Insert(v.ID); err != nil {
                c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Insert db models failed: " + err.Error()})
                return
            }
        }
        c.JSON(http.StatusOK, Resp{Code: 0, Data: "Init models success"})
        return
    }
    c.JSON(http.StatusOK, Resp{Code: 0, Data: "Init models success, No new models"})
}


/*
AllModels return all available models from db.

Path:

	/admin/model/all

Args:

	GET auth: Admin Token
*/
func AllModels(c *gin.Context) {
    modelSvc := svc.ModelSvc{}
    models, err := modelSvc.All()
    if err != nil {
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Read db models failed: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, Resp{
        Code: 0,
        Data: models,
    })
}