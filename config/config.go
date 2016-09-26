package config

import (
	"encoding/json"
	"github.com/pufferpanel/pufferd/logging"
	"io/ioutil"
)

var config map[string]interface{}

func Load() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		logging.Error("Error loading config", err)
		config = make(map[string]interface{})
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		logging.Error("Error loading config", err)
	}
}

func Get(key string) string {
	val := config[key]
	if val != nil {
		return val.(string)
	} else {
		return ""
	}
}

func GetOrDefault(key string, def string) string {
	val := Get(key)
	if val == "" {
		return def
	}
	return val
}