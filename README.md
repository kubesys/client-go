# client-go

We expect to provide a go client:
- **Flexibility**. It can support all Kubernetes-based systems with minimized extra development, such as [Openshift](https://www.redhat.com/en/technologies/cloud-computing/openshift), [istio](https://istio.io/), etc.
- **Usability**. Developers just need to learn to write json/yaml(kubernetes native style) from [Kubernetes documentation](https://kubernetes.io/docs/home/).
- **Integration**. It can work with the other Kubernetes clients, such as  [official](https://github.com/kubernetes-client/go).

## Comparison

|                           | [official](https://github.com/kubernetes-client/go) | [cdk8s](https://cdk8s.io/) | [this project](https://github.com/kubesys/kubernetes-client-go)  | 
|---------------------------|------------------|------------------|-------------------|
|        Compatibility                      | for kubernetes-native kinds    | for crd kinds                 |  for both |
|  Support customized Kubernetes resources  |  a lot of development          | a lot of development          |  zero-deployment     |
|    Works with the other SDKs              |  complex                       | complex                       |  simple              |     

## Architecture

![avatar](./docs/arch.png)
 
## Installation


```shell
git clone --recursive https://github.com/kubesys/client-go
```

### Maven users


## Usage

- [Usage](#usage)
    - [中文文档](https://www.yuque.com/kubesys/kubernetes-client/overview)
    - [Creating a client](#creating-a-client)
    - [Simple example](#simple-example)
    - [Get all kinds](#get-all-kinds)
    - [Work with other SDKs](#work-with-other-sdks)


### Creating a client


There are two ways to create a client:

- By url and token:

```go
client := new KubernetesClient(url, token);
client.Init()
```

Here, the token can be created and get by following commands:

1. create token

```yaml
kubectl create -f https://raw.githubusercontent.com/kubesys/client-go/master/account.yaml
```

2. get token

```kubectl
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep kubernetes-client | awk '{print $1}') | grep "token:" | awk -F":" '{print$2}' | sed 's/ //g'

```

- By kubeconfig:

```go
client, err := kubesys.NewKubernetesClientWithDefaultKubeConfig()
if err == nil {
    fmt.Println("Failed to get kubeconfig.")
    return
}
client.Init()
```


### simple-example

Assume you have a json:

```json
{
  "apiVersion": "v1",
  "kind": "Pod",
  "metadata": {
    "name": "busybox",
    "namespace": "default",
    "labels": {
      "test": "test"
    }
  }
}
```

List resources:

```go
client.ListResources("Pod")
```

Create a resource:

```go
client.CreateResource(json);
```

Get a resource:

```go
client.GetResource("Pod", "default", "busybox");
```

Delete a resource::

```go
client.DeleteResource("Pod", "default", "busybox")
```

### get-all-kinds

```go
fmt.Println(client.GetKinds());
```

## for developer

```
go mod init client-go
go mod tidy
```


## RoadMap

- 1.0.x: product ready
  - 1.0.0: using gojson
