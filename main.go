// Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
package main

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
)

func main() {

	url := "https://114.119.188.144:6443"
	tok := "eyJhbGciOiJSUzI1NiIsImtpZCI6IkZ5cEhSMzlCaEJvQkdNVHlxalNaTlFTeWNyb0NfajVKYTFlMkd0TFhrRVEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbi1zbHFkNiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6Ijk5M2RmZmU5LTExOGUtNDU1NS1iYzE3LWNmZjQ1YTRlMWJhOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.PEdgdWQa69fyJO-f5fTN1RsRXTmce2NGahyb7dCQMwulC0ObvoK6wOajZ22ZdS7cFfE7oP-gw49wTEFtL1iJWAUSxAN1idmdNusn77fLZl-a1njtq4DE26Bp7D1r0bharLqEEyCxW76UNTciXyVwyalrjB_Hn9e_lpOaUwJHyYOd0zB1_Fqj_2EOXQK4WBYpThDW3Pyf2UKuScwa1GoyFRsE4-Sc_Yosi4CPy34zP42nf970x-VfnLhBbHpygXVfMtC9PMV4795N5NQqR2GrrvdH6lD7gD-I9_LvSSxcQCAisv7qeuhVzXTmbeje5Hfd2Zh7DSYaYUMWpV1fz2BpPA"
	client := kubesys.NewKubernetesClient(url, tok)
	client.Init()

	//createResource(client)
	//getResource(client)
	//updateResource(client)
	//deleteResource(client)
	//listResources(client)

	//watchResources(client)
	watchResource(client)
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
	json, err := client.CreateResource(createPod())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(json.ToString())
}

func deleteResource(client *kubesys.KubernetesClient) {
	json, _ := client.DeleteResource("Pod", "default", "busybox")
	fmt.Println(json.ToString())
}

func getResource(client *kubesys.KubernetesClient) {
	json, _ := client.GetResource("Pod", "default", "busybox")
	fmt.Println(json.ToString())
}

func listResources(client *kubesys.KubernetesClient) {
	json,_ := client.ListResources("Deployment", "")
	fmt.Println(json.ToString())
}

func createPod() string {
	return "{\n  \"apiVersion\": \"v1\",\n  \"kind\": \"Pod\",\n  \"metadata\": {\n    \"name\": \"busybox\",\n    \"namespace\": \"default\"\n  },\n  \"spec\": {\n    \"containers\": [\n      {\n        \"image\": \"busybox\",\n        \"env\": [{\n           \"name\": \"abc\",\n           \"value\": \"abc\"\n        }],\n        \"command\": [\n          \"sleep\",\n          \"3600\"\n        ],\n        \"imagePullPolicy\": \"IfNotPresent\",\n        \"name\": \"busybox\"\n      }\n    ],\n    \"restartPolicy\": \"Always\"\n  }\n}"
}

func updateResource(client *kubesys.KubernetesClient) {

	labels := make(map[string]interface{})
	labels["test"] = "test"

	obj, _  := client.GetResource("Pod", "default", "busybox")
	metadata := obj.GetJsonObject("metadata")
	metadata.Put("labels", labels)
	fmt.Println(obj.ToString())

	json,err := client.UpdateResource(obj.ToString())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(json.ToString())
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
