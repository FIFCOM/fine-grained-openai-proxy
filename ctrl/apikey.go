package ctrl

import (
	"fine-grained-openai-proxy/svc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
AllApiKeys Get all OpenAI API keys

Path:

	/admin/apikey/all

Args:

	GET auth: Admin token
*/
func AllApiKeys(c *gin.Context) {
    apikeySvc := svc.ApiKeySvc{}
    apikeys, err := apikeySvc.All()
    if err != nil {
        c.JSON(http.StatusInternalServerError, Resp{Code: 1, Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, Resp{Code: 0, Data: apikeys})
}

/*
InsertApiKey Insert a new OpenAI API key

Path:

	/admin/apikey/insert

Args:

	GET auth: Admin token
	POST key: OpenAI API key
*/
func InsertApiKey(c *gin.Context) {
    key := c.PostForm("key")
    apikeySvc := svc.ApiKeySvc{}
    err := apikeySvc.Insert(key)
    if err != nil {
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Insert API Key Error: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, Resp{
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
func DeleteApiKey(c *gin.Context) {
    ids := c.PostForm("id")
    id, err := strconv.ParseInt(ids, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Delete API Key Error: " + err.Error()})
        return
    }

    apikeySvc := svc.ApiKeySvc{}
    err = apikeySvc.Delete(&svc.ApiKey{ID: id})
    if err != nil {
        c.JSON(http.StatusBadRequest, Resp{Code: 1, Error: "Delete API Key Error: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, Resp{
        Code: 0, Msg: "Delete API key successfully",
    })
}