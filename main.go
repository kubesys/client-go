// Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
package main

import (
	"./kubesys"
	"encoding/json"
	"fmt"
)

func main() {

	url := "https://119.8.188.235:6443"
	token := ""


	client := kubesys.NewKubernetesClient(url, token)
	client.Init()
	//fmt.Println(len(client.Analyzer.KindToFullKindMapper["Deployment"]))
	//fmt.Println(client.ListResources("Deployment", ""))
	//fmt.Println(client.GetResource("Pod", "default", "busybox"))
	//fmt.Println(client.DeleteResource("Pod", "default", "busybox"))
	//fmt.Println(client.CreateResource(createPod()))
	fmt.Println(client.UpdateResource(updatePod(client)))
	//watchResources(client)
	//watchResource(client)

	//json, _ := client.GetResource("Pod", "default", "busybox")
	//fmt.Println(json.GetObjectNode("metadata").GetString("name"))
}

func watchResource(client *kubesys.KubernetesClient) {
	watcher := kubesys.NewKubernetesWatcher(client)
	client.WatchResource("Pod", "default", "busybox", watcher)
}

func watchResources(client *kubesys.KubernetesClient) {
	watcher := kubesys.NewKubernetesWatcher(client)
	client.WatchResources("Pod", "", watcher)
}

func createPod() string {
	return "{\n  \"apiVersion\": \"v1\",\n  \"kind\": \"Pod\",\n  \"metadata\": {\n    \"name\": \"busybox\",\n    \"namespace\": \"default\"\n  },\n  \"spec\": {\n    \"containers\": [\n      {\n        \"image\": \"busybox\",\n        \"env\": [{\n           \"name\": \"abc\",\n           \"value\": \"abc\"\n        }],\n        \"command\": [\n          \"sleep\",\n          \"3600\"\n        ],\n        \"imagePullPolicy\": \"IfNotPresent\",\n        \"name\": \"busybox\"\n      }\n    ],\n    \"restartPolicy\": \"Always\"\n  }\n}"
}

func updatePod(client *kubesys.KubernetesClient) string {
	jsonObj, _  := client.GetResource("Pod", "default", "busybox")
	labels := make(map[string]interface{})
	labels["test"] = "test"
	jsonObj.GetMap("metadata")["labels"] = labels
	updateObj, _ := json.Marshal(jsonObj.Object)
	return string(updateObj)
}
