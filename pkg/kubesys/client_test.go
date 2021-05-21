package kubesys

import (
	"fmt"
	"testing"
)
var (
	url = ""
	token = ""
)

func TestClientWithLabelSelector(t *testing.T) {
	client := NewKubernetesClient(url, token)
	client.Init()
	filter := make(map[string]string)
	filter["app"] = "busybox1"
	pods, _ := client.ListResourcesWithLabelSelector("Pod", "", filter)
	fmt.Println(pods.GetArrayNode("items").Size())
}


func TestGetResourceNotExist(t *testing.T) {
	client := NewKubernetesClient(url, token)
	client.Init()

	pod1, err := client.GetResource("Pod", "default", "aaa")
	fmt.Println(err)
	pod2, err := client.GetResource("Pod", "default", "busybox1")
	fmt.Println(err)
	if pod1 != nil {
		t.Errorf("Expected nil, but get %s", pod1)
	}
	if pod2 == nil {
		t.Errorf("Expect %s, but get nil", "busybox1")
	}
}
