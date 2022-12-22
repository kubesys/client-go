package main

import (
	"github.com/kubesys/client-go/pkg/kubesys"
	"log"
	"time"
)

func main() {
	client := kubesys.NewKubernetesClientInCluster()
	client.Init()
	for {
		resp, err := client.ListResources("Pod", "in-cluster-ns")
		if err != nil {
			panic(err)
		}
		log.Println(1)
		log.Println(string(resp))
		time.Sleep(100 * time.Second)
	}

}
