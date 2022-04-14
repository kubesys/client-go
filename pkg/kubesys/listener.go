/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */

package kubesys

import (
	"fmt"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/4/8
 */
// TODO
func listen(client *KubernetesClient, registry *Registry) {

	crds, _ := client.ListResources("CustomResourceDefinition", "")

	items := ToJsonObject(crds).Get("items").Array()

	for i := 0; i < len(items); i++ {
		item := items[i]
		fmt.Println(item.String())
		group := item.Get("spec").Get("group").String()
		vers := item.Get("spec").Get("versions").Array()
		for j := 0; j < len(vers); j++ {
			ver := vers[i].Get("name").String()
			url := client.Url + "/apis/" + group + "/" + ver
			register(client, url, registry)
		}
	}
}
