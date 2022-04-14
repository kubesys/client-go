/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */

package kubesys

/**
 * this class is used for get Url for various kinds and operates in Kubernetes
 * see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/
 *
 *      author: wuheng@iscas.ac.cn
 *      date  : 2022/4/2
 *      since : v2.0.0
 */
func (client *KubernetesClient) baseUrl(fullKind string, namespace string) string {
	ruleBase := client.analyzer.RuleBase
	url := ruleBase.FullKindToApiPrefixMapper[fullKind] + "/"
	url += namespacePath(ruleBase.FullKindToNamespaceMapper[fullKind], namespace)
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

func (client *KubernetesClient) BindingResourceStatusUrl(fullKind string, namespace string, name string) string {
	return client.baseUrl(fullKind, namespace) + "/" + name + "/binding"
}
