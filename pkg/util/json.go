/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package util

import (
	"encoding/json"
	"fmt"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 */

type JsonNode struct {

}

type ObjectNode struct {
	JsonNode
	Object       map[string]interface{}
}

type ArrayNode struct {
	JsonNode
	Object       []interface{}
}

/************************************************************
 *
 *      initialization
 *
 *************************************************************/

func NewObjectNodeWithValue(value map[string]interface{}) *ObjectNode {
	json := new(ObjectNode)

	json.Object = value

	return json
}

func NewArrayNodeWithValue(value []interface{}) *ArrayNode {
	json := new(ArrayNode)

	json.Object = value

	return json
}

func (json *ObjectNode) GetObjectNode(key string) *ObjectNode {
	return NewObjectNodeWithValue(json.Object[key].(map[string]interface{}))
}

func (json *ArrayNode) GetObjectNode(idx int) *ObjectNode {
	return NewObjectNodeWithValue(json.Object[idx].(map[string]interface{}))
}

func (json *ObjectNode) GetArrayNode(key string) *ArrayNode {
	return NewArrayNodeWithValue(json.Object[key].([]interface{}))
}

func (json *ArrayNode) GetArrayNode(idx int) *ArrayNode {
	return NewArrayNodeWithValue(json.Object[idx].([]interface{}))
}

func (json *ObjectNode) GetMap(key string) map[string]interface{} {
	return json.Object[key].(map[string]interface{})
}

func (json *ArrayNode) GetMap(idx int) map[string]interface{} {
	return json.Object[idx].(map[string]interface{})
}

func (json *ObjectNode) GetArray(key string) []interface{} {
	return json.Object[key].([]interface{})
}

func (json *ArrayNode) GetArray(idx int) []interface{} {
	return json.Object[idx].([]interface{})
}

func (json *ObjectNode) GetString(key string) string {
	return json.Object[key].(string)
}


func (json *ArrayNode) GetString(idx int) string {
	return json.Object[idx].(string)
}

func (json *ObjectNode) GetBool(key string) bool {
	return json.Object[key].(bool)
}

func (json *ArrayNode) GetBool(idx int) bool {
	return json.Object[idx].(bool)
}

func (json *ObjectNode) GetInt(key string) int {
	return json.Object[key].(int)
}

func (json *ArrayNode) GetInt(idx int) int {
	return json.Object[idx].(int)
}

func (json *ArrayNode) Size() int {
	return len(json.Object)
}

func (obj *ObjectNode) Into(v interface{}) error {
	if obj == nil {
		return nil
	}
	objByte, err := json.Marshal(obj.Object)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(objByte, v)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}