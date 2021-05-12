package kubesys

import (
	"fmt"
	"testing"
)
var (
	url = "https://133.133.135.42:6443"
	token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjN2U3VZUk16R3ZfZGNaMkw4bVktVGlRWnJGZFB2NWprU1lrd0hObnNBVFEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbi10Z202ZyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjI0MjNlMDJmLTdmYzAtNDEzYi04ODczLTc0YTM3MTFkMzdkOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.KVJ7NC4NAWViLy2YkFFzzg0G4NcKnAZzw8VYooyXaLQlyfJWysR0giU8QLcSRs5BqIagff2EcVBuVHmSE4o1Zt3AMayStk-stwdtQre28adKYwR4aJLtfa1Wqmw--RiBHZmOjOmzynDdtWEe_sJPl4bGSxMvjFEKy6OepXOctnqZjUq4x2mMK-FID5hmeoHY6oAcfrRuAJsHRuLEAJQzLiMAf9heTuRNxcv3OTyfGtLOOj9risr59wilC_JWVPC5DC5TkEe4-8OeWg_mKA-lwSss_nyGMCsBqPIdPeyd3RQQ9ADPDq-JP2Nci0zoqOEwgZu3nQ3wOovR7lFBbRxsQQ"
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
