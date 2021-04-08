/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 */
type KubernetesClient struct {
	Url        string
	Token      string
	Http       *http.Client
	Analyzer   *KubernetesAnalyzer
}

/************************************************************
 *
 *      initialization
 *
 *************************************************************/
func NewKubernetesClient(url string, token string) *KubernetesClient {
	client := new(KubernetesClient)
	client.Url = url
	client.Token = token
	client.Http = &http.Client {Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	client.Analyzer = NewKubernetesAnalyzer()
	return client
}

func NewKubernetesClientWithAnalyzer(url string, token string, analyzer *KubernetesAnalyzer) *KubernetesClient {
	client := new(KubernetesClient)
	client.Url = url
	client.Token = token
	client.Http = &http.Client {Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	client.Analyzer = analyzer
	return client
}

func (client *KubernetesClient) Init() {
	client.Analyzer.Learning(*client)
}

/************************************************************
 *
 *      Common
 *
 *************************************************************/
func (client *KubernetesClient) RequestResource(request *http.Request) (map[string]interface{}, error) {
	res, err := client.Http.Do(request)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]interface{})
	json.Unmarshal(body, &result)
	return result, nil
}

func (client *KubernetesClient) CreateRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer " + client.Token)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, nil
}

func GetMapFromMap(values map[string]interface{}, key string) map[string]interface{} {
	return values[key].(map[string]interface{})
}

func GetArrayFromMap(values map[string]interface{}, key string) []interface{} {
	return values[key].([]interface{})
}

func getNamespace(supportNS bool, value string) string {
	if supportNS && len(value) != 0 {
		return "namespaces/" + value + "/"
	}
	return ""
}

/************************************************************
 *
 *      Core
 *
 *************************************************************/

func (client *KubernetesClient) CreateResource(jsonStr string) (*ObjectNode, error) {
	var jsonObj = make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &jsonObj)
	kind := jsonObj["kind"].(string)
	namespace := jsonObj["metadata"].(map[string]interface{})["namespace"].(string)
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind]
	req, _ := client.CreateRequest("POST", url, strings.NewReader(jsonStr))
	value, _ := client.RequestResource(req)
	return NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) UpdateResource(jsonStr string) (*ObjectNode, error) {
	fmt.Println(jsonStr)
	var jsonObj = make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &jsonObj)
	kind := jsonObj["kind"].(string)
	namespace := jsonObj["metadata"].(map[string]interface{})["namespace"].(string)
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind] + "/" + jsonObj["metadata"].(map[string]interface{})["name"].(string)
	req, _ := client.CreateRequest("PUT", url, strings.NewReader(jsonStr))
	value, _ := client.RequestResource(req)
	return NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) DeleteResource(kind string, namespace string, name string) (*ObjectNode, error) {
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind] + "/" + name
	req, _ := client.CreateRequest("DELETE", url, nil)
	value, _ := client.RequestResource(req)
	return NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) GetResource(kind string, namespace string, name string) (*ObjectNode, error) {
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind] + "/" + name
	req, _ := client.CreateRequest("GET", url, nil)
	value, _ := client.RequestResource(req)
	return NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) ListResources(kind string, namespace string) (*ObjectNode, error) {
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind]
	req, _ := client.CreateRequest("GET", url, nil)
	value, _ := client.RequestResource(req)
	return NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) WatchResource(kind string, namespace string, name string, watcher *KubernetesWatcher)  {
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/watch/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind] + "/" + name
	watcher.Watching(url)
}

func (client *KubernetesClient) WatchResources(kind string, namespace string, watcher *KubernetesWatcher)  {
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/watch/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind]
	watcher.Watching(url)
}