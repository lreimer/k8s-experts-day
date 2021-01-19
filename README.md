# Kubernetes Experts Day

Demo repository for the Kubernetes Experts Day session at ContainerConf 2020/21.

## Declarative Management of Kubernetes Objects Using Kustomize

```bash
# see https://kustomize.io
# see https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/

# preview of the Kustomize output
$ kubectl kustomize kustomized/base/
$ kubectl kustomize kustomized/overlays/int/

# apply the Kustomized output
$ kubectl apply -f kustomized/overlays/int/
$ kubectl delete -f kustomized/overlays/int/
```

## Using Kubernetes during Continuous Integration



## Using the Kubernetes API on the CLI

```bash
$ kubectl proxy
$ APISERVER=http://127.0.0.1:8001

$ curl -X GET $APISERVER/api
$ curl -X GET $APISERVER/apis

# see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/

# global watches
$ curl -X GET $APISERVER/api/v1/watch/events
$ curl -X GET $APISERVER/api/v1/watch/services
$ curl -X GET $APISERVER/api/v1/watch/pods

# namespace specific watches
$ curl -X GET $APISERVER/api/v1/watch/namespaces/default/events
$ curl -X GET $APISERVER/api/v1/watch/namespaces/default/pods

# deployment specific watches
$ curl -X GET $APISERVER/apis/apps/v1/
$ curl -X GET $APISERVER/apis/apps/v1/watch/deployments
$ curl -X GET $APISERVER/apis/apps/v1/namespaces/default/deployments\?watch\=true

$ curl -X GET $APISERVER/apis/batch/v1beta1/namespaces/{namespace}/cronjobs

# alternatively, directly against the API server
$ TOKEN=$(kubectl get secrets -o jsonpath="{.items[?(@.metadata.annotations['kubernetes\.io/service-account\.name']=='default')].data.token}"|base64 --decode)
$ APISERVER=$(kubectl config view -o jsonpath="{.clusters[?(@.name==\"$CLUSTER_NAME\")].cluster.server}")
$ curl -X GET $APISERVER/api --header "Authorization: Bearer $TOKEN" --insecure
```

## Using the Kubernetes API programmatically

```bash
# see https://kubernetes.io/docs/reference/using-api/client-libraries/
# see https://github.com/kubernetes/client-go
# see https://github.com/kubernetes-client/java

$ cd event-watcher-java/
$ ./gradlew clean ass
$ ./gradlew run
```

## Maintainer

M.-Leander Reimer (@lreimer), <mario-leander.reimer@qaware.de>

## License

This software is provided under the MIT open source license, read the `LICENSE` file for details.

