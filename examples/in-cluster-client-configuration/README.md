# Authenticating inside the cluster

This example shows you how to configure a client with client-go to authenticate
to the Kubernetes API from an application running inside the Kubernetes cluster.

client-go uses the [Service Account token][sa] mounted inside the Pod at the
`/var/run/secrets/kubernetes.io/serviceaccount` path when the
`kubesys.NewKubernetesClientInCluster()` is used.

## Running this example

First compile the application for Linux:

    cd in-cluster-client-configuration
    GOOS=linux go build -o ./app .

Then package it to a docker image using the provided Dockerfile to run it on
Kubernetes.


    docker build -t in-cluster .

If you have RBAC enabled on your cluster, use the following
snippet to create role binding which will grant the default service account view
permissions.

```
kubectl apply -f yamls/account.yaml
```

Notice that, in order to init the client, the pod need thr right to access the kubernetes address,
which is defined in ClusterRole' nonResourceURLs

```yaml
  - nonResourceURLs:
      - /
    verbs:
      - get
```

And you need to bind the account to system:discovery to get the access to other urls(/api, /apis, /api/*, /apis/*)

Then, run the image in a Pod with a single instance Deployment:

    kubectl apply -f yamls/deployment.yaml

The example now runs on Kubernetes API and successfully queries the number of
pods in the cluster every 100 seconds.

You can use `kubectl logs` to see the result. 

### Clean up

To stop this example and clean up the pod, press <kbd>Ctrl</kbd>+<kbd>C</kbd> on
the `kubectl run` command and then run:

    kubectl delete deployment demo

[sa]: https://kubernetes.io/docs/admin/authentication/#service-account-tokens
[mk]: https://kubernetes.io/docs/getting-started-guides/minikube/
