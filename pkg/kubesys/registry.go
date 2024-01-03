/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */

package kubesys

import (
	"encoding/json"
	"fmt"
)

/**
 * This class is used for
 *
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 *      since : v2.0.0
 */

type Registry struct {
	RuleBase *RuleBase
}

func NewRegistry(ruleBase *RuleBase) *Registry {
	registry := new(Registry)
	registry.RuleBase = ruleBase
	return registry
}

func register(client *KubernetesClient, url string, registry *Registry) {

	resourceRequest, err := client.createRequest("GET", url, nil)
	if err != nil {
		fmt.Println("createRequest error", err)
		panic(err)
	}

	resourceStringValues, err := client.doRequest(resourceRequest)

	if err != nil {
		fmt.Println("request resource string values error", err)
		panic(err)
	}
	resourceValues := make(map[string]interface{})
	json.Unmarshal([]byte(resourceStringValues), &resourceValues)

	apiVersion := resourceValues["groupVersion"].(string)
	for _, w := range resourceValues["resources"].([]interface{}) {
		resourceValue := w.(map[string]interface{})
		shortKind := resourceValue["kind"].(string)
		fullKind := getFullKind(resourceValue, shortKind, apiVersion)

		if _, ok := registry.RuleBase.FullKindToApiPrefixMapper[fullKind]; !ok {
			registry.RuleBase.KindToFullKindMapper[shortKind] = append(registry.RuleBase.KindToFullKindMapper[shortKind], fullKind)
			registry.RuleBase.FullKindToApiPrefixMapper[fullKind] = url

			registry.RuleBase.FullKindToNameMapper[fullKind] = resourceValue["name"].(string)
			registry.RuleBase.FullKindToNamespaceMapper[fullKind] = resourceValue["namespaced"].(bool)

			registry.RuleBase.FullKindToVersionMapper[fullKind] = apiVersion
			registry.RuleBase.FullKindToGroupMapper[fullKind] = getGroup(apiVersion)
			registry.RuleBase.FullKindToVerbsMapper[fullKind] = resourceValue["verbs"]
		}
	}
}
