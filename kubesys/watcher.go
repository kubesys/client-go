/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"bufio"
	"encoding/json"
	"fmt"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 */
type KubernetesWatcher struct {
	Client     *KubernetesClient
}

/************************************************************
 *
 *      initialization
 *
 *************************************************************/
func NewKubernetesWatcher(client *KubernetesClient) *KubernetesWatcher {
	watcher := new(KubernetesWatcher)

	watcher.Client = client

	return watcher
}

func (watcher *KubernetesWatcher) Watching(url string) {
	watcherClient := NewKubernetesClientWithAnalyzer(url, watcher.Client.Token, watcher.Client.Analyzer)
	req, _ := watcherClient.CreateRequest("GET", url, nil)
	resp, _ := watcherClient.Http.Do(req)
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		var jsonObj = make(map[string]interface{})
		json.Unmarshal([]byte(line), &jsonObj)

		if jsonObj["type"] == "ADDED" {
			watcher.Add(jsonObj["object"].(map[string]interface{}))
		} else if jsonObj["type"] == "MODIFIED" {
			watcher.Modify(jsonObj["object"].(map[string]interface{}))
		} else if jsonObj["type"] == "DELETED" {
			watcher.Delete(jsonObj["object"].(map[string]interface{}))
		}
	}
}

func (watcher *KubernetesWatcher) Add(obj map[string]interface{}) {
	fmt.Println("ADDED")
	fmt.Println(obj)
}

func (watcher *KubernetesWatcher) Modify(obj map[string]interface{}) {
	fmt.Println("MODIFIED")
	fmt.Println(obj)
}

func (watcher *KubernetesWatcher) Delete(obj map[string]interface{}) {
	fmt.Println("DELETED")
	fmt.Println(obj)
}
