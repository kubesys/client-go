/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	jsonObj "github.com/kubesys/kubernetes-client-go/pkg/json"
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
	Url      string
	Token    string
	Http     *http.Client
	Analyzer *KubernetesAnalyzer
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
	client.Http = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	client.Analyzer = NewKubernetesAnalyzer()
	return client
}

func NewKubernetesClientWithAnalyzer(url string, token string, analyzer *KubernetesAnalyzer) *KubernetesClient {
	client := new(KubernetesClient)
	client.Url = url
	client.Token = token
	client.Http = &http.Client{Transport: &http.Transport{
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
func (client *KubernetesClient) RequestResource(request *http.Request) (string, error) {
	res, err := client.Http.Do(request)
	if res.StatusCode != http.StatusOK {
		return "", errors.New("request status " + res.Status)
	}
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (client *KubernetesClient) CreateRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+client.Token)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, nil
}

func isNamespaced(supportNS bool, value string) string {
	if supportNS && len(value) != 0 {
		return "namespaces/" + value + "/"
	}
	return ""
}

// deprecated
func getRealKind(kind string, apiVersion string) string {
	index := strings.Index(apiVersion, "/")
	if index == -1 {
		return kind
	}
	return apiVersion[0:index] + "." + kind
}

func namespace(jsonObj *jsonObj.JsonObject) string {
	namespace := ""
	if jsonObj.GetJsonObject("metadata").HasKey("namespace")  {
		namespace, _ = jsonObj.GetJsonObject("metadata").GetString("namespace")
	}
	return namespace
}

func fullKind(jsonObj *jsonObj.JsonObject) string {
	kind, _ := jsonObj.GetString("kind")
	apiVersion, _ := jsonObj.GetString("apiVersion")
	index := strings.Index(apiVersion, "/")
	if index == -1 {
		return kind
	}
	return apiVersion[0:index] + "." + kind
}

func name(jsonObj *jsonObj.JsonObject) string {
	name, _ := jsonObj.GetJsonObject("metadata").GetString("name")
	return name
}

func (client *KubernetesClient) baseUrl(fullKind string, namespace string) string {
	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/"
	url += isNamespaced(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind]
	return url
}

func (client *KubernetesClient) getResponse(fullKind string, namespace string) string {
	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/"
	url += isNamespaced(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind]
	return url
}

func checkAndReturnRealKind(kind string, mapper map[string][]string) (string, error) {
	index := strings.Index(kind, ".")
	if index == -1 {
		if len(mapper[kind]) == 1 {
			return mapper[kind][0], nil
		} else if len(mapper[kind]) == 0 {
			return "", errors.New("invalid kind")
		} else {
			value := ""
			for _, s := range mapper[kind] {
				value += "," + s
			}
			return "", errors.New("please use fullKind: " + value[1:])
		}

	}
	return kind, nil
}


/************************************************************
 *
 *      Core
 *
 *************************************************************/

func (client *KubernetesClient) CreateResource(jsonStr string) (*jsonObj.JsonObject, error) {

	inputJson, err := jsonObj.ParseObject(jsonStr)
	if err != nil {
		return nil, err
	}

	url := client.baseUrl(fullKind(inputJson), namespace(inputJson))
	req, _ := client.CreateRequest("POST", url, strings.NewReader(jsonStr))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}

func (client *KubernetesClient) UpdateResource(jsonStr string) (*jsonObj.JsonObject, error) {

	inputJson, err := jsonObj.ParseObject(jsonStr)
	if err != nil {
		return nil, err
	}

	url := client.baseUrl(fullKind(inputJson), namespace(inputJson)) + "/" + name(inputJson)
	req, _ := client.CreateRequest("PUT", url, strings.NewReader(jsonStr))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}

func (client *KubernetesClient) DeleteResource(kind string, namespace string, name string) (*jsonObj.JsonObject, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.baseUrl(fullKind, namespace)  + "/" + name
	req, _ := client.CreateRequest("DELETE", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}

func (client *KubernetesClient) GetResource(kind string, namespace string, name string) (*jsonObj.JsonObject, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.baseUrl(fullKind, namespace)  + "/" + name
	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}

func (client *KubernetesClient) ListResources(kind string, namespace string) (*jsonObj.JsonObject, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.baseUrl(fullKind, namespace)
	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}

func (client *KubernetesClient) UpdateResourceStatus(jsonStr string) (*jsonObj.JsonObject, error) {
	inputJson, err := jsonObj.ParseObject(jsonStr)
	if err != nil {
		return nil, err
	}

	url := client.baseUrl(fullKind(inputJson), namespace(inputJson)) + "/" + name(inputJson) + "/status"
	req, _ := client.CreateRequest("PUT", url, strings.NewReader(jsonStr))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}

func (client *KubernetesClient) BindResources(pod *jsonObj.JsonObject, host string) (*jsonObj.JsonObject, error) {
	var podJson = make(map[string]interface{})
	podJson["apiVersion"] = "v1"
	podJson["kind"] = "Binding"

	var meta = make(map[string]interface{})
	meta["name"], _ = pod.GetJsonObject("metadata").GetString("name")
	meta["namespace"], _ = pod.GetJsonObject("metadata").GetString("namespace")
	podJson["metadata"] = meta

	var target = make(map[string]interface{})
	target["apiVersion"] = "v1"
	target["kind"] = "Node"
	target["name"] = host
	podJson["target"] = target

	fullKind := fullKind(pod)
	namespace := namespace(pod)
	url := client.baseUrl(fullKind, namespace) + "/" + name(pod) + "/binding"

	jsonBytes, _ := json.Marshal(podJson)
	req, _ := client.CreateRequest("POST", url, strings.NewReader(string(jsonBytes)))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}

func (client *KubernetesClient) WatchResource(kind string, namespace string, name string, watcher *KubernetesWatcher) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/watch/"
	url += isNamespaced(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind] + "/" + name
	watcher.Watching(url)
}

func (client *KubernetesClient) WatchResources(kind string, namespace string, watcher *KubernetesWatcher) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/watch/"
	url += isNamespaced(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind]
	watcher.Watching(url)
}

/************************************************************
 *
 *      Core for Object
 *
 *************************************************************/

func (client *KubernetesClient) CreateResourceObject(obj interface{}) (*jsonObj.JsonObject, error) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return client.CreateResource(string(jsonStr))
}

func (client *KubernetesClient) UpdateResourceObject(obj interface{}) (*jsonObj.JsonObject, error) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return client.UpdateResource(string(jsonStr))
}

/************************************************************
 *
 *      With Label Filter
 *
 *************************************************************/
func (client *KubernetesClient) ListResourcesWithLabelSelector(kind string, namespace string, labels map[string]string) (*jsonObj.JsonObject, error) {
	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.baseUrl(fullKind, namespace) + "?labelSelector="
	for key, value := range labels {
		url += key + "%3D" + value + ","
	}
	url = url[:len(url)-1]


	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
    if err != nil {
    	return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}

/************************************************************
 *
 *      With Field Filter
 *
 *************************************************************/
func (client *KubernetesClient) ListResourcesWithFieldSelector(kind string, namespace string, fields map[string]string) (*jsonObj.JsonObject, error) {
	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.baseUrl(fullKind, namespace) + "?fieldSelector="
	for key, value := range fields {
		url += key + "%3D" + value + ","
	}
	url = url[:len(url)-1]

	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	outputJson, err := jsonObj.ParseObject(value)
	if err != nil {
		return nil, err
	}

	return outputJson, nil
}
