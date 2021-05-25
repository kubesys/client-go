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

	ss := `{"kind":"Task","apiVersion":"doslab.io/v1","metadata":{"name":"task2","namespace":"default","uid":"b23a2359-9fd7-472e-9a2a-869152037861","resourceVersion":"8861789","generation":3,"creationTimestamp":"2021-05-18T08:13:15Z","annotations":{"gpu-limit":"0.2","gpu-memory":"1073741824","gpu-request":"0.2","kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"doslab.io/v1\",\"kind\":\"Task\",\"metadata\":{\"annotations\":{\"gpu-limit\":\"0.2\",\"gpu-memory\":\"1073741824\",\"gpu-request\":\"0.2\",\"speedup\":\"5.0\"},\"name\":\"task2\",\"namespace\":\"default\"},\"spec\":{\"containers\":[{\"command\":[\"sh\",\"-c\",\"nvidia-smi dmon -d 5\"],\"image\":\"pytorch/pytorch:1.6.0-cuda10.1-cudnn7-devel\",\"name\":\"torch\",\"volumeMounts\":[{\"mountPath\":\"/workspace\",\"name\":\"code\"}]}],\"volumes\":[{\"hostPath\":{\"path\":\"/root/gpushare/code\",\"type\":\"Directory\"},\"name\":\"code\"}]}}\n","schedule-gpuid":"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783","schedule-node":"dell04","schedule-time":"1621325603","speedup":"5.0"},"managedFields":[{"manager":"kubectl-client-side-apply","operation":"Update","apiVersion":"doslab.io/v1","time":"2021-05-18T08:13:15Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:gpu-limit":{},"f:gpu-memory":{},"f:gpu-request":{},"f:kubectl.kubernetes.io/last-applied-configuration":{},"f:speedup":{}}},"f:spec":{".":{},"f:volumes":{}}}},{"manager":"Go-http-client","operation":"Update","apiVersion":"doslab.io/v1","time":"2021-05-18T08:13:57Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{"f:schedule-gpuid":{},"f:schedule-node":{},"f:schedule-time":{}}},"f:spec":{"f:containers":{}},"f:status":{".":{},"f:pod_status":{".":{},"f:phase":{},"f:qosClass":{}}}}}]},"spec":{"volumes":[{"name":"code","hostPath":{"path":"/root/gpushare/code","type":"Directory"}}],"containers":[{"name":"torch","image":"pytorch/pytorch:1.6.0-cuda10.1-cudnn7-devel","command":["sh","-c","nvidia-smi dmon -d 5"],"resources":{},"volumeMounts":[{"name":"code","mountPath":"/workspace"}]}]},"status":{"pod_status":{"phase":"Running","conditions":[{"type":"Initialized","status":"True","lastProbeTime":"2021-05-18T08:39:08Z","lastTransitionTime":"2021-05-18T08:39:08Z"},{"type":"Ready","status":"False","lastProbeTime":"null","lastTransitionTime":"2021-05-18T08:43:48Z","reason":"ContainersNotReady","message":"containers with unready status: [torch]"},{"type":"ContainersReady","status":"False","lastProbeTime":"null","lastTransitionTime":"2021-05-18T08:43:48Z","reason":"ContainersNotReady","message":"containers with unready status: [torch]"},{"type":"PodScheduled","status":"True","lastProbeTime":"null","lastTransitionTime":"2021-05-18T08:39:08Z"}],"hostIP":"133.133.135.42","podIP":"192.168.228.109","podIPs":[{"ip":"192.168.228.109"}],"startTime":"2021-05-18T08:39:08Z","containerStatuses":[{"name":"torch","state":{"terminated":{"exitCode":137,"reason":"Error","startedAt":"2021-05-18T08:39:09Z","finishedAt":"2021-05-18T08:43:47Z","containerID":"docker://43872042b9a4bbe1f649fb65ad6097df7fbc7d50f940c0e738321b45bfda326e"}},"lastState":{},"ready":false,"restartCount":0,"image":"pytorch/pytorch:1.6.0-cuda10.1-cudnn7-devel","imageID":"docker-pullable://pytorch/pytorch@sha256:ccebb46f954b1d32a4700aaeae0e24bd68653f92c6f276a608bf592b660b63d7","containerID":"docker://43872042b9a4bbe1f649fb65ad6097df7fbc7d50f940c0e738321b45bfda326e","started":false}],"qosClass":"BestEffort"}}}`
	client.UpdateResource(ss)
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
