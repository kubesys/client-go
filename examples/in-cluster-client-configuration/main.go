package main

import (
	"github.com/kubesys/client-go/pkg/kubesys"
	"log"
	"time"
)

func main() {
	client := kubesys.NewKubernetesClientInCluster()
	for {
		resp, err := client.ListResources("Pod", "default")
		if err != nil {
			panic(err)
		}
		log.Println(resp)
		time.Sleep(100 * time.Second)
	}

}
