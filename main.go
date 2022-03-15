// Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
package main

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/client-go/pkg/kubesys"
)

func main() {

	url := ""
	tok := ""
	client := kubesys.NewKubernetesClient(url, tok)
	client.Init()

	// client, err := kubesys.NewKubernetesClientWithDefaultKubeConfig()
	// if err == nil {
	// 	fmt.Println("Failed to get kubeconfig.")
	// 	return
	// }
	// client.Init()

	//createResource(client)
	//getResource(client)
	//updateResource(client)
	//deleteResource(client)
	//listResources(client)

	//watchResources(client)
	//watchResource(client)
	//fmt.Println(client.GetKinds())
	//fmt.Println(client.GetFullKinds())
	fmt.Println(kubesys.ToJsonObject(client.GetKindDesc()).ToString())
}

func watchResource(client *kubesys.KubernetesClient) {
	watcher := kubesys.NewKubernetesWatcher(client, PrintWatchHandler{})
	client.WatchResource("Pod", "default", "busybox", watcher)
}

func watchResources(client *kubesys.KubernetesClient) {
	watcher := kubesys.NewKubernetesWatcher(client, PrintWatchHandler{})
	client.WatchResources("Pod", "", watcher)
}

func createResource(client *kubesys.KubernetesClient) {
	jsonRes, err := client.CreateResource(createPod())
	if err != nil {
		fmt.Println(err)
	}
	json := kubesys.ToJsonObject(jsonRes)
	fmt.Println(json.ToString())
}

func deleteResource(client *kubesys.KubernetesClient) {
	jsonRes, _ := client.DeleteResource("Pod", "default", "busybox")
	json := kubesys.ToJsonObject(jsonRes)
	fmt.Println(json.ToString())
}

func getResource(client *kubesys.KubernetesClient) {
	jsonRes, _ := client.GetResource("Pod", "default", "busybox")
	//fmt.Println(kubesys.ToJsonObject(jsonRes))
	fmt.Println(kubesys.ToGolangMap(jsonRes)["metadata"].(map[string]interface {})["name"].(string))
}

func listResources(client *kubesys.KubernetesClient) {
	jsonRes,_ := client.ListResources("Deployment", "")
	json := kubesys.ToJsonObject(jsonRes)
	fmt.Println(json.ToString())
}

func createPod() string {
	return "{\n  \"apiVersion\": \"v1\",\n  \"kind\": \"Pod\",\n  \"metadata\": {\n    \"name\": \"busybox\",\n    \"namespace\": \"default\"\n  },\n  \"spec\": {\n    \"containers\": [\n      {\n        \"image\": \"busybox\",\n        \"env\": [{\n           \"name\": \"abc\",\n           \"value\": \"abc\"\n        }],\n        \"command\": [\n          \"sleep\",\n          \"3600\"\n        ],\n        \"imagePullPolicy\": \"IfNotPresent\",\n        \"name\": \"busybox\"\n      }\n    ],\n    \"restartPolicy\": \"Always\"\n  }\n}"
}

func updateResource(client *kubesys.KubernetesClient) {

	labels := make(map[string]interface{})
	labels["test"] = "test"

	objRes, _  := client.GetResource("Pod", "default", "busybox")
	obj := kubesys.ToJsonObject(objRes)
	metadata := obj.GetJsonObject("metadata")
	metadata.Put("labels", labels)
	fmt.Println(metadata.ToString())

	obj.Put("metadata", metadata.ToInterface())
	fmt.Println(obj.ToString())

	jsonRes,err := client.UpdateResource(obj.ToString())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(kubesys.ToJsonObject(jsonRes).ToString())
}

type PrintWatchHandler struct {}

func (p PrintWatchHandler) DoAdded(obj map[string]interface{}) {
	json,_ :=json.Marshal(obj)
	fmt.Println("ADDED: " + string(json))
}
func (p PrintWatchHandler) DoModified(obj map[string]interface{}) {
	json,_ :=json.Marshal(obj)
	fmt.Println("MODIFIED: " + string(json))
}
func (p PrintWatchHandler) DoDeleted(obj map[string]interface{}) {
	json,_ :=json.Marshal(obj)
	fmt.Println("DELETED: " + string(json))
}

