/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"strings"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 */
func getGroup(apiVersion string) string {
	index := strings.LastIndex(apiVersion, "/")
	if index > 0 {
		return apiVersion[0:index]
	}
	return ""
}

func getFullKind(resourceValue map[string]interface{}, shortKind string, apiVersion string) string {
	index := strings.Index(apiVersion, "/")
	apiGroup := ""
	if index != -1 {
		apiGroup = apiVersion[0:index]
	}

	fullKind := ""
	if len(apiGroup) == 0 {
		fullKind = shortKind
	} else {
		fullKind = apiGroup + "." + shortKind
	}
	return fullKind
}
