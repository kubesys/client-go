// Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
package main

import (
	"./kubesys"
	"fmt"
)

func main() {

	url := "https://119.8.188.235:6443"
	token := ""

	client := kubesys.NewKubernetesClient(url, token)
	client.Init()
	fmt.Println(client.ListResources("Deployment", ""))
}
