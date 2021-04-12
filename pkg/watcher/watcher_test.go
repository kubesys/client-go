package watcher

import (
	"fmt"
	"testing"
	. "github.com/kubesys/kubernetes-client-go/pkg/client"
)

type PrintWatchHandler struct {}

func (p PrintWatchHandler) DoAdded(obj map[string]interface{}) {
	fmt.Println("add pod")
}
func (p PrintWatchHandler) DoModified(obj map[string]interface{}) {
	fmt.Println("update pod")
}
func (p PrintWatchHandler) DoDeleted(obj map[string]interface{}) {
	fmt.Println("delete pod")
}

func TestWatchHandler(t *testing.T) {
	url := ""
	token := ""
	client := NewKubernetesClient(url, token)
	client.Init()
	watcher := NewKubernetesWatcher(client, PrintWatchHandler{})
	client.WatchResource("Pod", "default", "busybox", watcher)
}
