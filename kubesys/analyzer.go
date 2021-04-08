// Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
package kubesys

//  author: wuheng@iscas.ac.cn
//  date: 2021/4/8
type KubernetesAnalyzer struct {
	KindToFullKindMapper              map[string]string
	FullKindToNameMapper              map[string]string
	FullKindToNamespaceMapper         map[string]string

	FullKindToVersionMapper           map[string]string
	FullKindToGroupMapper             map[string]string
	FullKindToVerbsMapper             map[string]string

	FullKindToApiPrefixMapper         map[string]string
}
