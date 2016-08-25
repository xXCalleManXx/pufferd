package config

import (
	"encoding/json"
	"io/ioutil"
)

var config map[string]interface{}

func Load() {
	data, _ := ioutil.ReadFile("config")
	json.Unmarshal(data, &config)
}

func Get(key string) string {
	val := config[key]
	if val != nil {
		return val.(string)
	} else {
		return ""
	}
}
