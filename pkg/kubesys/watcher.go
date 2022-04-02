/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package kubesys

import (
	"bufio"
	"encoding/json"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 */

type WatchHandler interface {
	DoAdded(obj map[string]interface{})
	DoModified(obj map[string]interface{})
	DoDeleted(obj map[string]interface{})
}

type KubernetesWatcher struct {
	Client  *KubernetesClient
	handler WatchHandler
}

/************************************************************
 *
 *      initialization
 *
 *************************************************************/

func NewKubernetesWatcher(client *KubernetesClient, handler WatchHandler) *KubernetesWatcher {
	return &KubernetesWatcher{
		Client:  client,
		handler: handler,
	}
}

func (watcher *KubernetesWatcher) Watching(url string) {
	watcherClient := NewKubernetesClientWithAnalyzer(url, watcher.Client.Token, watcher.Client.analyzer)
	watcherClient.http = watcher.Client.http
	req, _ := watcherClient.createRequest("GET", url, nil)
	resp, _ := watcherClient.http.Do(req)
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		var jsonObj = make(map[string]interface{})
		json.Unmarshal([]byte(line), &jsonObj)
		if jsonObj["type"] == "ADDED" {
			watcher.handler.DoAdded(jsonObj["object"].(map[string]interface{}))
		} else if jsonObj["type"] == "MODIFIED" {
			watcher.handler.DoModified(jsonObj["object"].(map[string]interface{}))
		} else if jsonObj["type"] == "DELETED" {
			watcher.handler.DoDeleted(jsonObj["object"].(map[string]interface{}))
		}
	}
}
