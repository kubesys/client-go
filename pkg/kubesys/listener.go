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
func listen(client KubernetesClient, registry *Registry) {

	crds,_ := client.ListResources("CustomResourceDefinition", "")

	items := ToJsonObject(crds).GetJsonArray("items")

	for i := 0; i < len(items.Values()); i++ {
		item := items.GetJsonObject(i)
		group, _ := item.GetJsonObject("spec").GetString("group")
		vers := item.GetJsonArray("versions")
		for j := 0; j < len(vers.Values()); j++ {
			ver, _ := vers.GetJsonObject(j).GetString("name")
			fmt.Println("/apis/" + group + "/" + ver)
		}
	}
}


