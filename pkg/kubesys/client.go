/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	jsonObj "github.com/kubesys/client-go/pkg/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
	if strings.HasSuffix(url, "/") {
		client.Url = url[0 : len(url)-1]
	} else {
		client.Url = url
	}
	client.Token = token
	client.Http = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	client.Analyzer = NewKubernetesAnalyzer()
	return client
}

func NewKubernetesClientWithDefaultKubeConfig() (*KubernetesClient, error) {
	client, err := NewKubernetesClientWithKubeConfig("/etc/kubernetes/admin.conf")
	if err == nil {
		return client, err
	}
	return NewKubernetesClientWithKubeConfig(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
}

func NewKubernetesClientWithKubeConfig(kubeConfig string) (*KubernetesClient, error) {
	config, err := NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	client := new(KubernetesClient)
	client.Url = config.Server
	httpClient, err := HTTPClientFor(config)
	if err != nil {
		return nil, err
	}
	client.Http = httpClient
	client.Analyzer = NewKubernetesAnalyzer()
	return client, nil
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
	client.Analyzer.Learning(client)
}

/************************************************************
 *
 *      Common
 *
 *************************************************************/
func (client *KubernetesClient) RequestResource(request *http.Request) ([]byte, error) {
	res, err := client.Http.Do(request)
	if res.StatusCode != http.StatusOK {
		if err != nil {
			return nil, errors.New("request status " + res.Status + ": " + err.Error())
		} else {
			return nil, errors.New("request status " + res.Status)
		}
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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

func namespace(jsonObj *jsonObj.JsonObject) string {
	namespace := ""
	if jsonObj.GetJsonObject("metadata").HasKey("namespace") {
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

func kind(fullKind string) string {
	index := strings.LastIndex(fullKind, ".")
	if index == -1 {
		return fullKind
	}
	return fullKind[index+1:]
}

func name(jsonObj *jsonObj.JsonObject) string {
	name, _ := jsonObj.GetJsonObject("metadata").GetString("name")
	return name
}

func (client *KubernetesClient) getResponse(fullKind string, namespace string) string {
	ruleBase := client.Analyzer.RuleBase
	url := ruleBase.FullKindToApiPrefixMapper[fullKind] + "/"
	url += isNamespaced(ruleBase.FullKindToNamespaceMapper[fullKind], namespace)
	url += ruleBase.FullKindToNameMapper[fullKind]
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

func (client *KubernetesClient) CreateResource(jsonStr string) ([]byte, error) {

	inputJson, err := jsonObj.ParseObject(jsonStr)
	if err != nil {
		return nil, err
	}

	url := client.CreateResourceUrl(fullKind(inputJson), namespace(inputJson))
	req, _ := client.CreateRequest("POST", url, strings.NewReader(jsonStr))
	_, err = client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return []byte(jsonStr), nil
}

func (client *KubernetesClient) UpdateResource(jsonStr string) ([]byte, error) {

	inputJson, err := jsonObj.ParseObject(jsonStr)
	if err != nil {
		return nil, err
	}

	url := client.UpdateResourceUrl(fullKind(inputJson), namespace(inputJson), name(inputJson))
	req, _ := client.CreateRequest("PUT", url, strings.NewReader(jsonStr))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) DeleteResource(kind string, namespace string, name string) ([]byte, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.DeleteResourceUrl(fullKind, namespace, name)
	req, _ := client.CreateRequest("DELETE", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) GetResource(kind string, namespace string, name string) ([]byte, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.GetResourceUrl(fullKind, namespace, name)
	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) ListResources(kind string, namespace string) ([]byte, error) {

	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.ListResourcesUrl(fullKind, namespace)
	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) UpdateResourceStatus(jsonStr string) ([]byte, error) {
	inputJson, err := jsonObj.ParseObject(jsonStr)
	if err != nil {
		return nil, err
	}

	url := client.UpdateResourceStatusUrl(fullKind(inputJson), namespace(inputJson), name(inputJson))
	req, _ := client.CreateRequest("PUT", url, strings.NewReader(jsonStr))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) BindResources(pod *jsonObj.JsonObject, host string) ([]byte, error) {
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
	url := client.BindingResourceStatusUrl(fullKind, namespace, name(pod))

	jsonBytes, _ := json.Marshal(podJson)
	req, _ := client.CreateRequest("POST", url, strings.NewReader(string(jsonBytes)))
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) WatchResource(kind string, namespace string, name string, watcher *KubernetesWatcher) {

	ruleBase := client.Analyzer.RuleBase
	fullKind, err := checkAndReturnRealKind(kind, ruleBase.KindToFullKindMapper)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := ruleBase.FullKindToApiPrefixMapper[fullKind] + "/watch/"
	url += isNamespaced(ruleBase.FullKindToNamespaceMapper[fullKind], namespace)
	url += ruleBase.FullKindToNameMapper[fullKind] + "/" + name
	url += "/?watch=true&timeoutSeconds=315360000"
	watcher.Watching(url)
}

func (client *KubernetesClient) WatchResources(kind string, namespace string, watcher *KubernetesWatcher) {

	ruleBase := client.Analyzer.RuleBase
	fullKind, err := checkAndReturnRealKind(kind, ruleBase.KindToFullKindMapper)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := ruleBase.FullKindToApiPrefixMapper[fullKind] + "/watch/"
	url += isNamespaced(ruleBase.FullKindToNamespaceMapper[fullKind], namespace)
	url += ruleBase.FullKindToNameMapper[fullKind]
	url += "/?watch=true&timeoutSeconds=315360000"
	watcher.Watching(url)
}

/************************************************************
 *
 *      With Label Filter
 *
 *************************************************************/
func (client *KubernetesClient) ListResourcesWithLabelSelector(kind string, namespace string, labels map[string]string) ([]byte, error) {
	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.ListResourcesUrl(fullKind, namespace) + "?labelSelector="
	for key, value := range labels {
		url += key + "%3D" + value + ","
	}
	url = url[:len(url)-1]

	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

/************************************************************
 *
 *      With Field Filter
 *
 *************************************************************/
func (client *KubernetesClient) ListResourcesWithFieldSelector(kind string, namespace string, fields map[string]string) ([]byte, error) {
	fullKind, err := checkAndReturnRealKind(kind, client.Analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.ListResourcesUrl(fullKind, namespace) + "?fieldSelector="
	for key, value := range fields {
		url += key + "%3D" + value + ","
	}
	url = url[:len(url)-1]

	req, _ := client.CreateRequest("GET", url, nil)
	value, err := client.RequestResource(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

/************************************************************
 *
 *      Metadata
 *
 *************************************************************/

func (client *KubernetesClient) GetKinds() []string {
	i := 0
	mapper := client.Analyzer.RuleBase.KindToFullKindMapper
	array := make([]string, len(mapper))
	for key, _ := range mapper {
		array[i] = key
		i++
	}
	return array
}

func (client *KubernetesClient) GetFullKinds() []string {
	i := 0
	mapper := client.Analyzer.RuleBase.FullKindToNameMapper
	array := make([]string, len(mapper))
	for key, _ := range mapper {
		array[i] = key
		i++
	}
	return array
}

func (client *KubernetesClient) GetKindDesc() []byte {
	var desc = make(map[string]interface{})

	ruleBase := client.Analyzer.RuleBase
	for fullKind, _ := range ruleBase.FullKindToNameMapper {
		var value = make(map[string]interface{})
		value["apiVersion"] = ruleBase.FullKindToVersionMapper[fullKind]
		value["kind"] = kind(fullKind)
		value["plural"] = ruleBase.FullKindToNameMapper[fullKind]
		value["verbs"] = ruleBase.FullKindToVerbsMapper[fullKind]
		desc[fullKind] = value
	}
	bytes, _ := json.Marshal(desc)
	return bytes
}
