package utils

func GetStringOrDefault(data map[string]interface{}, key string, def *string) string {
	if data == nil {
		return *def
	}
	var section = data[key]
	if section == nil {
		return *def
	} else {
		return section.(string)
	}
}

func GetBooleanOrDefault(data map[string]interface{}, key string, def bool) bool {
	if data == nil {
		return def
	}
	var section = data[key]
	if section == nil {
		return def
	} else {
		return section.(bool)
	}
}

func GetMapOrNull(data map[string]interface{}, key string) map[string]interface{} {
	if data == nil {
		return (map[string]interface{})(nil)
	}
	var section = data[key]
	if section == nil {
		return (map[string]interface{})(nil)
	} else {
		return section.(map[string]interface{})
	}
}

func GetObjectArrayOrNull(data map[string]interface{}, key string) []interface{} {
	if data == nil {
		return ([]interface{})(nil)
	}
	var section = data[key]
	if section == nil {
		return ([]interface{})(nil)
	} else {
		return section.([]interface{})
	}
}

func GetStringArrayOrNull(data map[string]interface{}, key string) []string {
	if data == nil {
		return ([]string)(nil)
	}
	var section = data[key]
	if section == nil {
		return ([]string)(nil)
	} else {
		var sec = section.([]interface{})
		var newArr = make([]string, len(sec))
		for i := 0; i < len(sec); i++ {
			newArr[i] = sec[i].(string)
		}
		return newArr
	}
}
