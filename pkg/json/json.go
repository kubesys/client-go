/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package json

import "encoding/json"

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/9/29
 */
func ParseObject(jsonStr string) (*JsonObject, error) {
	value := new(Value)

	err := json.Unmarshal([]byte(jsonStr), &value.data)
	if err != nil {
		return nil, err
	}

	return value.JsonObject(), nil
}

func ParseArray(jsonStr string) (*JsonArray, error) {
	value := new(Value)

	err := json.Unmarshal([]byte(jsonStr), &value.data)
	if err != nil {
		return nil, err
	}
	return value.JsonArray(), nil
}