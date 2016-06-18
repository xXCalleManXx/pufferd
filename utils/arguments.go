/*
 Copyright 2016 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package utils

import (
	"strings"
)

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
