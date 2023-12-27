/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */

package kubesys

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

/**
 * This class is used for get all apis from kube-apiserver analysis
 * see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/
 *
 *      author: wuheng@iscas.ac.cn
 *      date  : 2022/4/3
 *      since : v2.0.0
 */
func extract(client *KubernetesClient, registry *Registry) {
	// request kube-apiserver, such as http://IP:6443.
	registryRequest, err := client.createRequest("GET", client.Url, nil)
	if err != nil {
		errors.New(err.Error())
	}

	registryStringValues, err := client.doRequest(registryRequest)
	if err != nil {
		fmt.Println("request registry string values error, ", err)
		panic(err)
	}
	// if it is successful, the output is.
	// {
	//    "paths": [
	//        "/.well-known/openid-configuration",
	//        "/api",
	//        "/api/v1",
	//        "/apis",
	//        "/apis/",
	//        "/apis/admissionregistration.k8s.io",
	//        "/apis/admissionregistration.k8s.io/v1",
	//        "/apis/apiextensions.k8s.io",
	//        "/apis/apiextensions.k8s.io/v1",
	//        "/apis/apiregistration.k8s.io",
	//        "/apis/apiregistration.k8s.io/v1",
	//        "/apis/apps",
	//        "/apis/apps/v1",
	//        "/apis/authentication.k8s.io",
	//        "/apis/authentication.k8s.io/v1",
	//        "/apis/authorization.k8s.io",
	//        "/apis/authorization.k8s.io/v1",
	//        "/apis/autoscaling",
	//        "/apis/autoscaling/v1",
	//        "/apis/autoscaling/v2",
	//        "/apis/autoscaling/v2beta1",
	//        "/apis/autoscaling/v2beta2",
	//        "/apis/batch",
	//        "/apis/batch/v1",
	//        "/apis/batch/v1beta1",
	//        "/apis/certificates.k8s.io",
	//        "/apis/certificates.k8s.io/v1",
	//        "/apis/coordination.k8s.io",
	//        "/apis/coordination.k8s.io/v1",
	//        "/apis/discovery.k8s.io",
	//        "/apis/discovery.k8s.io/v1",
	//        "/apis/discovery.k8s.io/v1beta1",
	//        "/apis/events.k8s.io",
	//        "/apis/events.k8s.io/v1",
	//        "/apis/events.k8s.io/v1beta1",
	//        "/apis/flowcontrol.apiserver.k8s.io",
	//        "/apis/flowcontrol.apiserver.k8s.io/v1beta1",
	//        "/apis/flowcontrol.apiserver.k8s.io/v1beta2",
	//        "/apis/kubeovn.io",
	//        "/apis/kubeovn.io/v1",
	//        "/apis/networking.k8s.io",
	//        "/apis/networking.k8s.io/v1",
	//        "/apis/node.k8s.io",
	//        "/apis/node.k8s.io/v1",
	//        "/apis/node.k8s.io/v1beta1",
	//        "/apis/policy",
	//        "/apis/policy/v1",
	//        "/apis/policy/v1beta1",
	//        "/apis/rbac.authorization.k8s.io",
	//        "/apis/rbac.authorization.k8s.io/v1",
	//        "/apis/scheduling.k8s.io",
	//        "/apis/scheduling.k8s.io/v1",
	//        "/apis/storage.k8s.io",
	//        "/apis/storage.k8s.io/v1",
	//        "/apis/storage.k8s.io/v1beta1",
	//        "/version"
	//    ]
	registryValues := make(map[string]interface{})
	json.Unmarshal([]byte(registryStringValues), &registryValues)

	for _, v := range registryValues["paths"].([]interface{}) {
		path := v.(string)
		// just check /api and /apis
		if strings.HasPrefix(path, "/api") &&
			// go to /apis/node.k8s.io/v1 rather than /apis/node.k8s.io, or goto /api/v1
			(len(strings.Split(path, "/")) == 4 || strings.EqualFold(path, "/api/v1")) {
			register(client, client.Url+path, registry)
		}
	}
}
