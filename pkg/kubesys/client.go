/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */

package kubesys

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
 * this class is used for creating a connection between users' application and Kubernetes server.
 * It provides an easy-to-use way to Create, Update, Delete, Get, List and Watch all Kubernetes resources.
 *
 *      author: wuheng@iscas.ac.cn
 *      date  : 2022/4/1
 *      since : v2.0.0
 */

/************************************************************
 *
 *      struct
 *
 *************************************************************/

type KubernetesClient struct {
	Url      string              // required, user input
	Token    string              // required, user input
	http     *http.Client        // required, automatically created based on Url and Token
	analyzer *KubernetesAnalyzer // required, user input or automatically register all Kubernetes resources based on Http
}

/************************************************************
 *
 *      initialization
 *
 *************************************************************/

func createClient(url string, token string, http *http.Client, analyzer *KubernetesAnalyzer) *KubernetesClient {
	// init a NewKubernetesClient object
	client := new(KubernetesClient)

	// assignment
	client.Url = checkedUrl(url)
	client.Token = checkedToken(token)
	client.http = http
	client.analyzer = analyzer

	// return
	return client
}

func NewKubernetesClient(url string, token string) *KubernetesClient {
	return createClient(url, token, &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true},
		}}, NewKubernetesAnalyzer())
}

// NewKubernetesClientWithDefaultKubeConfig TODO
func NewKubernetesClientWithDefaultKubeConfig() *KubernetesClient {
	// filepath.Join(os.Getenv("HOME"), ".kube", "config")
	return NewKubernetesClientWithKubeConfig("/etc/kubernetes/admin.conf")
}

// NewKubernetesClientWithKubeConfig TODO
func NewKubernetesClientWithKubeConfig(kubeConfig string) *KubernetesClient {
	config, err := NewForConfig(kubeConfig)
	if err != nil {
		panic(err)
	}

	httpClient, err := HTTPClientFor(config)
	if err != nil {
		panic(err)
	}

	return createClient(config.Server, "", httpClient, NewKubernetesAnalyzer())
}

func NewKubernetesClientWithAnalyzer(url string, token string, analyzer *KubernetesAnalyzer) *KubernetesClient {
	return createClient(url, token, &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true},
		}}, analyzer)
}

func (client *KubernetesClient) Init() {
	// not initialized
	if len(client.analyzer.RuleBase.KindToFullKindMapper) == 0 {
		// initialing
		client.analyzer.Learning(client)
	}
}

/************************************************************
 *
 *      Core
 *
 *************************************************************/

func (client *KubernetesClient) CreateResource(jsonStr string) ([]byte, error) {

	inputJson := gjson.Parse(jsonStr)

	url := client.CreateResourceUrl(fullKind(inputJson), namespace(inputJson))

	req, err := client.createRequest("POST", url, strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}

	_, err = client.doRequest(req)
	if err != nil {
		return nil, err
	}

	return []byte(jsonStr), nil
}

func (client *KubernetesClient) UpdateResource(jsonStr string) ([]byte, error) {

	inputJson := gjson.Parse(jsonStr)

	url := client.UpdateResourceUrl(fullKind(inputJson), namespace(inputJson), name(inputJson))
	req, err := client.createRequest("PUT", url, strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}

	value, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) DeleteResource(kind string, namespace string, name string) ([]byte, error) {

	fullKind, err := toFullKind(kind, client.analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.DeleteResourceUrl(fullKind, namespace, name)
	req, err := client.createRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	value, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) GetResource(kind string, namespace string, name string) ([]byte, error) {

	fullKind, err := toFullKind(kind, client.analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.GetResourceUrl(fullKind, namespace, name)
	req, err := client.createRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	value, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) ListResources(kind string, namespace string) ([]byte, error) {

	fullKind, err := toFullKind(kind, client.analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.ListResourcesUrl(fullKind, namespace)
	req, err := client.createRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	value, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) UpdateResourceStatus(jsonStr string) ([]byte, error) {
	inputJson := gjson.Parse(jsonStr)

	url := client.UpdateResourceStatusUrl(fullKind(inputJson), namespace(inputJson), name(inputJson))
	req, err := client.createRequest("PUT", url, strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}

	value, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// BindResources TODO
func (client *KubernetesClient) BindResources(pod gjson.Result, host string) ([]byte, error) {
	var podJson = make(map[string]interface{})
	podJson["apiVersion"] = "v1"
	podJson["kind"] = "Binding"

	var meta = make(map[string]interface{})
	meta["name"] = pod.Get("metadata").Get("name").String()
	meta["namespace"] = pod.Get("metadata").Get("namespace").String()
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
	req, _ := client.createRequest("POST", url, strings.NewReader(string(jsonBytes)))
	value, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (client *KubernetesClient) WatchResource(kind string, namespace string, name string, watcher *KubernetesWatcher) {

	ruleBase := client.analyzer.RuleBase
	fullKind, err := toFullKind(kind, ruleBase.KindToFullKindMapper)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := ruleBase.FullKindToApiPrefixMapper[fullKind] + "/watch/"
	url += namespacePath(ruleBase.FullKindToNamespaceMapper[fullKind], namespace)
	url += ruleBase.FullKindToNameMapper[fullKind] + "/" + name
	url += "/?watch=true&timeoutSeconds=315360000"
	watcher.Watching(url)
}

func (client *KubernetesClient) WatchResources(kind string, namespace string, watcher *KubernetesWatcher) {

	ruleBase := client.analyzer.RuleBase
	fullKind, err := toFullKind(kind, ruleBase.KindToFullKindMapper)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := ruleBase.FullKindToApiPrefixMapper[fullKind] + "/watch/"
	url += namespacePath(ruleBase.FullKindToNamespaceMapper[fullKind], namespace)
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
	fullKind, err := toFullKind(kind, client.analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.ListResourcesUrl(fullKind, namespace) + "?labelSelector="
	for key, value := range labels {
		url += key + "%3D" + value + ","
	}
	url = url[:len(url)-1]

	req, _ := client.createRequest("GET", url, nil)
	value, err := client.doRequest(req)
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
	fullKind, err := toFullKind(kind, client.analyzer.RuleBase.KindToFullKindMapper)
	if err != nil {
		return nil, err
	}

	url := client.ListResourcesUrl(fullKind, namespace) + "?fieldSelector="
	for key, value := range fields {
		url += key + "%3D" + value + ","
	}
	url = url[:len(url)-1]

	req, _ := client.createRequest("GET", url, nil)
	value, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	return value, nil
}

/************************************************************
 *
 *      Common
 *
 *************************************************************/

func (client *KubernetesClient) doRequest(request *http.Request) ([]byte, error) {
	res, err := client.http.Do(request)
	if err != nil {
		return nil, errors.New("request error:" + err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("wrong request status: " + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *KubernetesClient) createRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	if len(client.Token) != 0 {
		req.Header.Add("Authorization", "Bearer "+client.Token)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func name(jsonObj gjson.Result) string {
	return jsonObj.Get("metadata").Get("name").String()
}

func namespace(jsonObj gjson.Result) string {

	if jsonObj.Get("metadata").Get("namespace").Exists() {
		return jsonObj.Get("metadata").Get("namespace").String()
	}
	return ""
}

// after command 'kubectl api-resources', you can see
// NAME                              SHORTNAMES   APIVERSION                             NAMESPACED   KIND
// pods                              po           v1                                     true         Pod
// deployments                       deploy       apps/v1                                true         Deployment
func kind(fullKind string) string {
	index := strings.LastIndex(fullKind, ".")
	if index == -1 {
		// for pods, the kind equals fullKind
		return fullKind
	}
	// for deployments, the fullKind is apps.Deployment, the kind is Deployment
	return fullKind[index+1:]
}

// after command 'kubectl api-resources', you can see
// NAME                              SHORTNAMES   APIVERSION                             NAMESPACED   KIND
// pods                              po           v1                                     true         Pod
// deployments                       deploy       apps/v1                                true         Deployment
func fullKind(jsonObj gjson.Result) string {
	kind := jsonObj.Get("kind").String()
	//apiVersion := jsonObj.Get("kind").Get("apiVersion").String()
	apiVersion := jsonObj.Get("apiVersion").String()

	index := strings.Index(apiVersion, "/")
	if index == -1 {
		// for pods, the fullKind equals kind
		return kind
	}
	// for deployments, the fullKind is apps.Deployment
	return apiVersion[0:index] + "." + kind
}

func namespacePath(supportNS bool, ns string) string {
	if supportNS && len(ns) != 0 {
		// if a kind supports namespace, and namespace is not null
		return "namespaces/" + ns + "/"
	}
	return ""
}

func toFullKind(kind string, mapper map[string][]string) (string, error) {
	index := strings.Index(kind, ".")
	if index == -1 {
		// it is just kind, we need to get fullKind
		if len(mapper[kind]) == 0 {
			return "", errors.New("wrong kind, please invoking 'GetKinds'")
		} else if len(mapper[kind]) == 1 {
			return mapper[kind][0], nil
		} else {
			// multiple fullKinds have a same kind
			value := ""
			for _, s := range mapper[kind] {
				value += "," + s
			}
			return "", errors.New("please input fullKind: " + value[1:])
		}

	}
	return kind, nil
}

/************************************************************
 *
 *      Metadata
 *
 *************************************************************/

func (client *KubernetesClient) GetKinds() []string {

	mapper := client.analyzer.RuleBase.KindToFullKindMapper

	i := 0
	keys := make([]string, len(mapper))
	for key, _ := range mapper {
		keys[i] = key
		i++
	}

	return keys
}

func (client *KubernetesClient) GetFullKinds() []string {

	mapper := client.analyzer.RuleBase.FullKindToNameMapper

	i := 0
	keys := make([]string, len(mapper))
	for key, _ := range mapper {
		keys[i] = key
		i++
	}
	return keys
}

func (client *KubernetesClient) GetKindDesc() []byte {

	desc := make(map[string]interface{})
	ruleBase := client.analyzer.RuleBase

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
