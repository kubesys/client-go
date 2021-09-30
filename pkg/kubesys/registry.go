/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 */
type Registry struct {
	RuleBase     *RuleBase
}

func NewRegistry(ruleBase *RuleBase) *Registry {
	registry := new(Registry)
	registry.RuleBase = ruleBase
	return registry
}

func Register(apiPrefix string, ruleBase *RuleBase, resourceValue map[string]interface{}, apiVersion string) {

	shortKind := resourceValue["kind"].(string)
	fullKind := getFullKind(resourceValue, shortKind, apiVersion)

	if _, ok := ruleBase.FullKindToApiPrefixMapper[fullKind]; !ok {
		ruleBase.KindToFullKindMapper[shortKind] = append(ruleBase.KindToFullKindMapper[shortKind], fullKind)
		ruleBase.FullKindToApiPrefixMapper[fullKind] = apiPrefix

		ruleBase.FullKindToNameMapper[fullKind] = resourceValue["name"].(string)
		ruleBase.FullKindToNamespaceMapper[fullKind] = resourceValue["namespaced"].(bool)

		ruleBase.FullKindToVersionMapper[fullKind] = apiVersion
		ruleBase.FullKindToGroupMapper[fullKind] = getGroup(apiVersion)
		ruleBase.FullKindToVerbsMapper[fullKind] = resourceValue["verbs"]
	}
}


