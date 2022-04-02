/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"regexp"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/9/30
 */

func checkedUrl(url string) string {
	// just support https without suffix '/'
	// this validation is not restricted
	httpsRegExp := regexp.MustCompile("https:\\/\\/([\\w.]+\\/?)\\S*")
	if !httpsRegExp.MatchString(url) {
		panic("just support https without suffix '/'")
	}
	return url
}

func checkedToken(token string) string {
	// just support token in Kubernetes Secret
	if len(token) != 950 {
		panic("just support token in Kubernetes Secret")
	}
	return token
}

func ToJsonObject(bytes []byte) gjson.Result {
	return gjson.Parse(string(bytes))
}

func ToGolangMap(bytes []byte) map[string]interface{} {
	values := make(map[string]interface{})
	json.Unmarshal([]byte(bytes), &values)
	return values
}
