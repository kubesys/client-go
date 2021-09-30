/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

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

