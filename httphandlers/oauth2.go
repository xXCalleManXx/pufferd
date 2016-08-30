package httphandlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/config"
	"github.com/pufferpanel/pufferd/logging"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func OAuth2Handler(gin *gin.Context) {
	authHeader := gin.Request.Header.Get("Authorization")
	if authHeader == "" {
		gin.AbortWithStatus(401)
		return
	}
	authArr := strings.SplitN(authHeader, " ", 2)
	if len(authArr) < 2 || authArr[0] != "Bearer" {
		gin.AbortWithStatus(401)
		return
	}
	ParseToken(authArr[1], gin)
}

func ParseToken(accessToken string, gin *gin.Context) {
	authUrl := config.Get("authserver")
	token := config.Get("authtoken")
	client := &http.Client{}
	data := url.Values{}
	data.Set("token", accessToken)
	request, _ := http.NewRequest("POST", authUrl, bytes.NewBufferString(data.Encode()))
	request.Header.Add("Authorization", "Bearer" + token)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(request)
	if err != nil {
		logging.Error("Error talking to auth server", err)
		gin.AbortWithStatus(500)
		return
	}
	if response.StatusCode != 200 {
		logging.Error("Error talking to auth server", response.StatusCode)
		gin.AbortWithStatus(500)
		return
	}
	var respArr map[string]interface{}
	json.NewDecoder(response.Body).Decode(&respArr)
	if respArr["error"] != nil {
		gin.AbortWithStatus(500)
		return
	}
	if respArr["active"].(bool) == false {
		gin.AbortWithStatus(401)
		return
	}
	gin.Set("server_id", respArr["server_id"].(string))
	gin.Set("scopes", strings.Split(respArr["scope"].(string), " "))
}
