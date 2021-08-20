/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/util"
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
func (client *KubernetesClient) RequestResource(request *http.Request) (map[string]interface{}, error) {
	res, err := client.Http.Do(request)
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("request status " + res.Status)
	}
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

	req.Header.Add("Authorization", "Bearer "+client.Token)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, nil
}

func getNamespace(supportNS bool, value string) string {
	if supportNS && len(value) != 0 {
		return "namespaces/" + value + "/"
	}
	return ""
}

func getRealKind(kind string, apiVersion string) string {
	index := strings.Index(apiVersion, "/")
	if index == -1 {
		return kind
	}
	return apiVersion[0:index] + "." + kind
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

func (client *KubernetesClient) CreateResource(jsonStr string) (*util.ObjectNode, error) {
	var jsonObj = make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &jsonObj)
	kind := getRealKind(jsonObj["kind"].(string), jsonObj["apiVersion"].(string))
	namespace := ""
	if _, ok := jsonObj["metadata"].(map[string]interface{})["namespace"]; ok {
		namespace = jsonObj["metadata"].(map[string]interface{})["namespace"].(string)
	}
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind]
	req, _ := client.CreateRequest("POST", url, strings.NewReader(jsonStr))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}
	return util.NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) UpdateResource(jsonStr string) (*util.ObjectNode, error) {
	var jsonObj = make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &jsonObj)
	kind := getRealKind(jsonObj["kind"].(string), jsonObj["apiVersion"].(string))
	namespace := ""
	if _, ok := jsonObj["metadata"].(map[string]interface{})["namespace"]; ok {
		namespace = jsonObj["metadata"].(map[string]interface{})["namespace"].(string)
	}
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind] + "/" + jsonObj["metadata"].(map[string]interface{})["name"].(string)
	req, _ := client.CreateRequest("PUT", url, strings.NewReader(jsonStr))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}
	return util.NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) DeleteResource(kind string, namespace string, name string) (*util.ObjectNode, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)

	if err != nil {
		return nil, err
	}

	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind] + "/" + name
	req, _ := client.CreateRequest("DELETE", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}
	return util.NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) GetResource(kind string, namespace string, name string) (*util.ObjectNode, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)

	if err != nil {
		return nil, err
	}

	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind] + "/" + name
	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}
	return util.NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) ListResources(kind string, namespace string) (*util.ObjectNode, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)

	if err != nil {
		return nil, err
	}

	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind]
	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}
	return util.NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) BindResources(pod *util.ObjectNode, host string) (*util.ObjectNode, error) {
	var jsonObj = make(map[string]interface{})
	jsonObj["apiVersion"] = "v1"
	jsonObj["kind"] = "Binding"

	var meta = make(map[string]interface{})
	meta["name"] = pod.GetObjectNode("metadata").GetString("name")
	meta["namespace"] = pod.GetObjectNode("metadata").GetString("namespace")
	jsonObj["metadata"] = meta

	var target = make(map[string]interface{})
	target["apiVersion"] = "v1"
	target["kind"] = "Node"
	target["name"] = host
	jsonObj["target"] = target

	kind := getRealKind(pod.GetString("kind"), pod.GetString("apiVersion"))
	namespace := pod.GetObjectNode("metadata").GetString("namespace")
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind] + "/"
	url += pod.GetObjectNode("metadata").GetString("name") + "/binding"

	jsonBytes, _ := json.Marshal(jsonObj)
	req, _ := client.CreateRequest("POST", url, strings.NewReader(string(jsonBytes)))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}
	return util.NewObjectNodeWithValue(value), nil
}

func (client *KubernetesClient) WatchResource(kind string, namespace string, name string, watcher *KubernetesWatcher) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/watch/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
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
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind]
	watcher.Watching(url)
}

func (client *KubernetesClient) UpdateResourceStatus(jsonStr string) (*util.ObjectNode, error) {
	var jsonObj = make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &jsonObj)
	kind := getRealKind(jsonObj["kind"].(string), jsonObj["apiVersion"].(string))
	namespace := ""
	if _, ok := jsonObj["metadata"].(map[string]interface{})["namespace"]; ok {
		namespace = jsonObj["metadata"].(map[string]interface{})["namespace"].(string)
	}
	url := client.Analyzer.FullKindToApiPrefixMapper[kind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[kind], namespace)
	url += client.Analyzer.FullKindToNameMapper[kind] + "/" + jsonObj["metadata"].(map[string]interface{})["name"].(string)
	url += "/status"
	req, _ := client.CreateRequest("PUT", url, strings.NewReader(jsonStr))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}
	return util.NewObjectNodeWithValue(value), nil
}

/************************************************************
 *
 *      Core for Object
 *
 *************************************************************/

func (client *KubernetesClient) CreateResourceObject(obj interface{}) (*util.ObjectNode, error) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return client.CreateResource(string(jsonStr))
}

func (client *KubernetesClient) UpdateResourceObject(obj interface{}) (*util.ObjectNode, error) {
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
func (client *KubernetesClient) ListResourcesWithLabelSelector(kind string, namespace string, labels map[string]string) (*util.ObjectNode, error) {
	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.KindToFullKindMapper)

	if err != nil {
		return nil, err
	}

	url := client.Analyzer.FullKindToApiPrefixMapper[fullKind] + "/"
	url += getNamespace(client.Analyzer.FullKindToNamespaceMapper[fullKind], namespace)
	url += client.Analyzer.FullKindToNameMapper[fullKind]
	url += "?labelSelector="
	for key, value := range labels {
		url += key + "%3D" + value
	}
	req, _ := client.CreateRequest("GET", url, nil)
	value, _ := client.RequestResource(req)
	return util.NewObjectNodeWithValue(value), nil
}
