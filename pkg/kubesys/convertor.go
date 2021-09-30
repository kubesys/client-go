/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 */
func (client *KubernetesClient) baseUrl(fullKind string, namespace string) string {
	ruleBase := client.Analyzer.RuleBase
	url := ruleBase.FullKindToApiPrefixMapper[fullKind] + "/"
	url += isNamespaced(ruleBase.FullKindToNamespaceMapper[fullKind], namespace)
	url += ruleBase.FullKindToNameMapper[fullKind]
	return url
}

func (client *KubernetesClient) CreateResourceUrl(fullKind string, namespace string) string {
	return client.baseUrl(fullKind, namespace)
}

func (client *KubernetesClient) ListResourcesUrl(fullKind string, namespace string) string {
	return client.baseUrl(fullKind, namespace)
}

func (client *KubernetesClient) UpdateResourceUrl(fullKind string, namespace string, name string) string {
	return client.baseUrl(fullKind, namespace) + "/" + name
}

func (client *KubernetesClient) DeleteResourceUrl(fullKind string, namespace string, name string) string {
	return client.baseUrl(fullKind, namespace) + "/" + name
}

func (client *KubernetesClient) GetResourceUrl(fullKind string, namespace string, name string) string {
	return client.baseUrl(fullKind, namespace) + "/" + name
}

func (client *KubernetesClient) UpdateResourceStatusUrl(fullKind string, namespace string, name string) string {
	return client.baseUrl(fullKind, namespace) + "/" + name + "/status"
}



