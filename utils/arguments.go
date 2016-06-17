package utils

import "strings"

func ReplaceTokens(msg string, mapping map[string]string) string {
	newmsg := msg
	for key, value := range mapping {
		newmsg = strings.Replace(newmsg, "${"+key+"}", value, -1)
	}
	return newmsg
}

func ReplaceTokensInArr(msg []string, mapping map[string]string) []string {
	newarr := make([]string, len(msg))
	for index, element := range msg {
		newarr[index] = ReplaceTokens(element, mapping)
	}
	return newarr
}
