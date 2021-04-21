# 多云平台服务通信及管理构件

现有的Kubernetes客户端(Java frabic8)存在支持参数有限、灵活性不够、不支持KVM虚拟机、DSL定义复杂等问题，难以支持多云平台服务的扩展。本构件分析了Kubernetes、KVM和第三方云平台的API特征，是一款新的定制化的多云平台服务管理客户端。

## 技术架构
如图1所示，本管理构件的核心是学习器，通过分析和学习Kubernetes原生资源（Pod）和用户自定义资源进行分析，根据URL规则，形成知识。<br>当用户根据API进行资源访问时，其核心思想就是根据类型进行查询，映射为Kubernetes的URL，进行代理执行。

![图1 多云平台服务通信及管理构件架构图](https://raw.githubusercontent.com/kubesys/kubernetes-client-go/main/pics/arch.jpg)
## 技术特色
如何面向Kubernetes的自定义资源类型提供管理客户端仍面临巨大挑战，已有的客户端框架或适配工作量巨大（如Java fabric8），学习成本高，或主要采用代码生成技术（code-gen），灵活性不够。此外，Kubernetes版本迭代频繁且存在不兼容问题，可能需要频繁更新开发框架版本。<br>为应对上述问题，本文设计一款学习驱动的，基于JSON的Kubernetes客户端，通过学习自修正以应对Kubernetes版本迭代频繁且不兼容问题，采用JSON解决已有框架学习成本高和灵活性不足的问题。

## 代码结构
主要代码结构如下：
```
.
├── go.mod
├── main.go
└── pkg
    ├── kubesys
    │   ├── analyzer.go
    │   ├── client.go
    │   ├── watcher.go
    │   └── watcher_test.go
    └── util
        └── json.go

3 directories, 9 files
```
pkg/kubesys包下包含本构件的主要代码<br>
analyzer.go 主要功能为对Kubernetes原生资源与自定义资源进行分析，根据URL规则，形成知识库<br>
client.go 主要功能为创建客户端，对http请求进行封装，为常用的创建、更新、删除资源等操作提供接口<br>
watcher.go 主要功能为监听资源变化，可以用来实现自定义的资源的控制逻辑<br>
<br>
util/包下包含本构件的相关工具
json.go 主要功能为对Golang提供的原生JSON包进行封装，抽象出JsonNode、ObjectNode和ArrayNode，提供更加方便的JSON处理接口。

## 部署方式
使用本构件的方式为，在go语言项目的go.mod中引入本项目作为依赖即可。
```
module test
go 1.15
require (
	github.com/kubesys/kubernetes-client-go v0.0.0-20210412025431-29031f2cac5f
)
```
## 使用说明
* 获取Kubernetes集群访问API Server的URL和token 
```
kubectl create -f https://raw.githubusercontent.com/kubesys/kubernetes-client-go/main/account.yaml

kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep kubernetes-client | awk '{print $1}') | grep "token:" | awk -F":" '{print$2}' | sed 's/ //g'

```
* 创建客户端
```go
url := "https://xxx.xxx.xxx.xxx:6443"
token := "<your token>"
client := kubesys.NewKubernetesClient(url, token)
client.Init()
```

* 创建资源
```
jsonStr := `{
  "apiVersion": "v1",
  "kind": "Pod",
  "metadata": {
    "name": "busybox",
    "namespace": "default",
    "labels": {
      "test": "test"
    }
  }
}`
client.CreateResource(jsonStr)
```
* 查询资源
```
client.GetResource("Pod", "default", "busybox")
```
* 更新资源
```
client.UpdateResource()
```
* 删除资源
```
client.DeleteResource("Pod", "default", "busybox")
```
