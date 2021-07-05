// Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
package main

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/client-go/pkg/kubesys"
)

func main() {

	url := "https://39.106.40.190:6443"
	tok := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImVDUzZnUkJ2OHI0VVA2VWpkdU1SLWNXTlF4aEhLQjNyamU2ZHhsd014cWMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbi16ZGxyZiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImJlZGJjMGUzLWJjNDAtNGQ2Zi1iMTAxLTk3ODkzOGZjYTZhNyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.Qmfmf4QwubartSyLJqIW2gXHdlyKlqQsNknIVtRNjfaydw6qCa8XuGS6egqPwiN-Al8GaoGuVflyJy-bolj-aVWY-a-9fWUB0itV4SdYTNeQV5hYv6sbhnuvSo3nHp2jyZjlRyvEQNxKyQaJF6eodJPjgzCVoj8BhsSqTu7vbzCTEMEnIz8AMGJLF9G6JuffBTpO83Ch_hVbquQnKQJjK60911D-5S6SD3SilQyk_WdYblorxbRXsSm8VNkHz6BWrfa7uCDcw46XnfVVuCyRKOGmIAeIWIDq2uaI6nECkcujWCCwYzEePq-SsXU4MRwFAYd-Rdt9Q8JUw9njR0I-5A"
	client := kubesys.NewKubernetesClient(url, tok)
	client.Init()
	obj1,err1 := client.CreateResource(kubesys.Issue2())
	fmt.Println(err1)
	fmt.Println(obj1)
	obj2,err2 := client.CreateResource(kubesys.Issue2())
	fmt.Println(err2)
	fmt.Println(obj2)
	//client.GetResource("Pod", "default", "busybox")
	//fmt.Println(len(client.Analyzer.KindToFullKindMapper["Deployment"]))
	//fmt.Println(client.ListResources("Deployment", ""))
	//fmt.Println(client.GetResource("Pod", "default", "busybox"))
	//fmt.Println(client.DeleteResource("Pod", "default", "busybox"))
	//fmt.Println(client.CreateResource(createPod()))
	//fmt.Println(client.UpdateResource(updatePod(client)))
	//watchResources(client)
	//watchResource(client)
	//json, _ := client.GetResource("Pod", "default", "busybox")
	//fmt.Println(json.GetObjectNode("metadata").GetString("name"))
}

//func watchResource(client *KubernetesClient) {
//	watcher := NewKubernetesWatcher(client)
//	client.WatchResource("Pod", "default", "busybox", watcher)
//}
//
//func watchResources(client *KubernetesClient) {
//	watcher := NewKubernetesWatcher(client)
//	client.WatchResources("Pod", "", watcher)
//}

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
