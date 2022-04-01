/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/9/30
 */

func ToJsonObject(bytes []byte) gjson.Result {
	return gjson.Parse(string(bytes))
}

func ToGolangMap(bytes []byte) map[string]interface{} {
	values := make(map[string]interface{})
	json.Unmarshal([]byte(bytes), &values)
	return values
}
