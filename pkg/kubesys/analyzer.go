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
	RuleBase     *RuleBase
}

func NewKubernetesAnalyzer() *KubernetesAnalyzer {
	ruleBase := new(RuleBase)

	ruleBase.KindToFullKindMapper = make(map[string][]string)
	ruleBase.FullKindToApiPrefixMapper = make(map[string]string)

	ruleBase.FullKindToNameMapper = make(map[string]string)
	ruleBase.FullKindToNamespaceMapper = make(map[string]bool)

	ruleBase.FullKindToVersionMapper = make(map[string]string)
	ruleBase.FullKindToGroupMapper = make(map[string]string)
	ruleBase.FullKindToVerbsMapper = make(map[string]interface{})

	analyzer := new(KubernetesAnalyzer)
	analyzer.RuleBase = ruleBase

	return analyzer
}

func (analyzer *KubernetesAnalyzer) Learning(client KubernetesClient) {
	extract(client, analyzer.RuleBase)
}

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
