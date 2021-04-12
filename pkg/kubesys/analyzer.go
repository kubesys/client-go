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
type KubernetesAnalyzer struct {
	KindToFullKindMapper              map[string][]string
	FullKindToApiPrefixMapper         map[string]string

	FullKindToNameMapper              map[string]string
	FullKindToNamespaceMapper         map[string]bool

	FullKindToVersionMapper           map[string]string
	FullKindToGroupMapper             map[string]string
	FullKindToVerbsMapper             map[string]interface{}
}

func NewKubernetesAnalyzer() *KubernetesAnalyzer {
	analyzer := new(KubernetesAnalyzer)

	analyzer.KindToFullKindMapper      = make(map[string][]string)
	analyzer.FullKindToApiPrefixMapper = make(map[string]string)

	analyzer.FullKindToNameMapper      = make(map[string]string)
	analyzer.FullKindToNamespaceMapper = make(map[string]bool)

	analyzer.FullKindToVersionMapper   = make(map[string]string)
	analyzer.FullKindToGroupMapper     = make(map[string]string)
	analyzer.FullKindToVerbsMapper     = make(map[string]interface{})

	return analyzer
}

func (analyzer *KubernetesAnalyzer) Learning (client KubernetesClient) {
	registryRequest,_ := client.CreateRequest("GET", client.Url, nil)
	registryValues, _ := client.RequestResource(registryRequest)
	for _, v := range registryValues["paths"].([]interface{}) {
		path := v.(string)
		if strings.HasPrefix(path, "/api") &&
			(len(strings.Split(path, "/")) == 4 ||
				strings.EqualFold(path, "/api/v1")) {
			resourceRequest,_ := client.CreateRequest("GET", client.Url + path, nil)
			resourceValues, _ := client.RequestResource(resourceRequest)
			apiVersion := resourceValues["groupVersion"].(string)
			for _, w := range resourceValues["resources"].([]interface{}) {
				resourceValue := w.(map[string]interface{})

				shortKind  := resourceValue["kind"].(string)
				fullKind := getFullKind(resourceValue, shortKind, apiVersion)

				if _, ok := analyzer.FullKindToApiPrefixMapper[fullKind]; !ok {
					analyzer.KindToFullKindMapper[shortKind] = append(analyzer.KindToFullKindMapper[shortKind], fullKind)
					analyzer.FullKindToApiPrefixMapper[fullKind] = client.Url + path

					analyzer.FullKindToNameMapper[fullKind] = resourceValue["name"].(string)
					analyzer.FullKindToNamespaceMapper[fullKind] = resourceValue["namespaced"].(bool)

					analyzer.FullKindToVersionMapper[fullKind] = apiVersion
					analyzer.FullKindToGroupMapper[fullKind] = getGroup(client.Url + path)
					analyzer.FullKindToVerbsMapper[fullKind] = resourceValue["verbs"]
				}
			}
		}
	}
}

func getGroup(url string) string {
	if strings.HasSuffix(url, "/api/v1") {
		return ""
	}
	stx := strings.LastIndex(url, "/")
	etx := strings.LastIndex(url[0:stx], "/")
	if etx > stx {
		return url[stx+1: etx]
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