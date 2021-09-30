/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"encoding/json"
	"strings"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/9/30
 */
type RuleBase struct {
	KindToFullKindMapper      map[string][]string
	FullKindToApiPrefixMapper map[string]string

	FullKindToNameMapper      map[string]string
	FullKindToNamespaceMapper map[string]bool

	FullKindToVersionMapper map[string]string
	FullKindToGroupMapper   map[string]string
	FullKindToVerbsMapper   map[string]interface{}
}

func NewKubernetesAnalyzer() *RuleBase {
	analyzer := new(RuleBase)

	analyzer.KindToFullKindMapper = make(map[string][]string)
	analyzer.FullKindToApiPrefixMapper = make(map[string]string)

	analyzer.FullKindToNameMapper = make(map[string]string)
	analyzer.FullKindToNamespaceMapper = make(map[string]bool)

	analyzer.FullKindToVersionMapper = make(map[string]string)
	analyzer.FullKindToGroupMapper = make(map[string]string)
	analyzer.FullKindToVerbsMapper = make(map[string]interface{})

	return analyzer
}

func (analyzer *RuleBase) Learning(client KubernetesClient) {
	registryRequest, _ := client.CreateRequest("GET", client.Url, nil)
	registryStringValues, _ := client.RequestResource(registryRequest)

	registryValues := make(map[string]interface{})
	json.Unmarshal([]byte(registryStringValues), &registryValues)

	for _, v := range registryValues["paths"].([]interface{}) {
		path := v.(string)
		if strings.HasPrefix(path, "/api") && (len(strings.Split(path, "/")) == 4 || strings.EqualFold(path, "/api/v1")) {
			resourceRequest, _ := client.CreateRequest("GET", client.Url+path, nil)
			resourceStringValues, _ := client.RequestResource(resourceRequest)

			resourceValues := make(map[string]interface{})
			json.Unmarshal([]byte(resourceStringValues), &resourceValues)

			apiVersion := resourceValues["groupVersion"].(string)
			for _, w := range resourceValues["resources"].([]interface{}) {
				resourceValue := w.(map[string]interface{})

				shortKind := resourceValue["kind"].(string)
				fullKind := getFullKind(resourceValue, shortKind, apiVersion)

				if _, ok := analyzer.FullKindToApiPrefixMapper[fullKind]; !ok {
					analyzer.KindToFullKindMapper[shortKind] = append(analyzer.KindToFullKindMapper[shortKind], fullKind)
					analyzer.FullKindToApiPrefixMapper[fullKind] = client.Url + path

					analyzer.FullKindToNameMapper[fullKind] = resourceValue["name"].(string)
					analyzer.FullKindToNamespaceMapper[fullKind] = resourceValue["namespaced"].(bool)

					analyzer.FullKindToVersionMapper[fullKind] = apiVersion
					analyzer.FullKindToGroupMapper[fullKind] = getGroup(apiVersion)
					analyzer.FullKindToVerbsMapper[fullKind] = resourceValue["verbs"]
				}
			}
		}
	}
}
