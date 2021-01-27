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
$ kubectl apply -k kustomized/overlays/int/
$ kubectl delete -k kustomized/overlays/int/
```

## Imperative Management of Kubernetes Objects Using Pulumi

```bash
# see https://www.pulumi.com/docs/get-started/kubernetes/
$ brew install pulumi
$ pulumi plugin install resource kubernetes v2.5.1

# using TypeScript as language
$ mkdir -p pulumi-demo-ts && cd pulumi-demo-ts
$ pulumi new kubernetes-typescript

$ pulumi up
$ pulumi destroy

$ cp ../nginx-deployment.yaml .
$ kube2pulumi typescript -f nginx-deployment.yaml

# using Go as language
$ mkdir -p pulumi-demo-go && cd pulumi-demo-go
$ pulumi new kubernetes-go

$ pulumi up
$ pulumi destroy

$ cp ../nginx-deployment.yaml .
$ kube2pulumi go -f nginx-deployment.yaml
```

## Using Kubernetes for CI/CD

```bash
# Continuous Load Testing with K6 on Kubernetes
# see https://github.com/lreimer/continuous-k6k8s

# next you can deploy the K6 stack with InfluxDB and Grafana
$ kubectl apply -f continuous-k6k8s.yaml

# open Grafana and import on of these K6 load test dashboards
# - see https://grafana.com/dashboards/2587
# - see https://grafana.com/grafana/dashboards/4411
$ open http://localhost:3000

# run adhoc tests as a simple pod
# be sure to pass the --restart flag, otherwise the containers gets restarted
$ kubectl run k6-nginx-test --image lreimer/k6-nginx-test --restart=Never --attach
$ kubectl delete pod/k6-nginx-test

# Continuous Security Tests with ZAP on Kubernetes
# https://github.com/lreimer/continuous-zapk8s

# Cloud Native CI/CD with Tekton
# see https://tekton.dev
$ kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
$ kubectl apply -f https://github.com/tektoncd/dashboard/releases/latest/download/tekton-dashboard-release.yaml
$ kubectl get pods -n tekton-pipelines

$ kubectl -n tekton-pipelines port-forward svc/tekton-dashboard 9097:9097
$ open http://localhost:9097

$ kubectl apply -f tekton-demos/task-hello.yaml
$ kubectl tkn task start hello
$ kubectl tkn taskrun logs --last -f 

$ kubectl apply -f task-goodbye.yaml
$ kubectl tkn task start goodbye 
$ kubectl tkn taskrun logs --last -f 

$ kubectl apply -f tekton-demos/pipeline-hello-goodbye.yaml
$ kubectl tkn pipeline start hello-goodbye
$ kubectl tkn pipelinerun logs --last -f 

# use Tekton triggers to run pipelines
# see https://tekton.dev/docs/triggers/install/
# see https://github.com/tektoncd/triggers/tree/v0.9.1/docs/getting-started
```

## Using the Kubernetes API on the CLI

```bash
# see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/

$ kubectl get events
$ kubectl get events -o json

$ kubectl get events --field-selector type!=Normal
$ kubectl get events --field-selector type=Normal,involvedObject.kind=Pod

$ kubectl proxy
$ APISERVER=http://127.0.0.1:8001

$ curl -X GET $APISERVER/api
$ curl -X GET $APISERVER/apis

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

## Defining and Using Custom Resources in Kubernetes

```bash
# see https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/
# see https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/

$ kubectl apply -f custom-resources/microservice-crd.yaml
$ kubectl get crds
$ kubectl apply -f custom-resources/microservice-example.yaml
$ kubectl get yservice

$ kubectl describe yservice microservice-example
$ curl -X GET $APISERVER/apis/k8s.qaware.de/v1alpha1/watch/microservices

$ kubectl apply -f custom-resources/supersecret-crd.yaml
$ kubectl get crds
$ kubectl apply -f custom-resources/supersecret-example.yaml
$ kubectl get supersecrets

$ kubectl describe supersecret supersecret-example
$ curl -X GET $APISERVER/apis/k8s.qaware.de/v1alpha1/watch/supersecrets
```

## Developing Custom Controllers with the Operator SDK

```bash
# see https://github.com/lreimer/graal-operators
# see https://github.com/lreimer/go-for-operations

# create operator project using SDK
$ mkdir -p operator-demo && cd operator-demo
$ operator-sdk init --project-version="2" --domain qaware.de --license none --owner "Mario-Leander Reimer" --plugins go.kubebuilder.io/v2 --repo github.com/lreimer/k8s-experts-day/operator-demo
$ operator-sdk create api --group k8s --version v1 --kind Demo --resource=true --controller=true

# build and install CRD
$ make install
$ kubectl get crds
$ kubectl describe crd demoes.k8s.qaware.de  

# first run
$ make run ENABLE_WEBHOOKS=false
$ kubectl apply -f config/samples/k8s_v1_demo.yaml
$ kubectl delete -f config/samples/k8s_v1_demo.yaml

# edit api/v1/demo_types.go to modify CRD
$ make generate manifests
$ make install
$ kubectl describe crd demoes.k8s.qaware.de  

# download test harness and perform release
# see https://sdk.operatorframework.io/docs/building-operators/golang/references/envtest-setup/
$ curl https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/master/hack/setup-envtest.sh -o setup-envtest.sh
$ chmod +x setup-envtest.sh

$ make docker-build docker-push
$ make deploy 
$ kubectl get all -n operator-demo-system
$ make undeploy
```

## Developing Operators with Kudo

```bash
$ kubectl krew install kudo
$ kubectl kudo init --unsafe-self-signed-webhook-ca

$ kubectl kudo install ./nginx-kudoperator
$ kubectl kudo get instances
```

## Developing Admission Controllers

![](https://d33wubrfki0l68.cloudfront.net/af21ecd38ec67b3d81c1b762221b4ac777fcf02d/7c60e/images/blog/2019-03-21-a-guide-to-kubernetes-admission-controllers/admission-controller-phases.png)

```bash
# see https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/
# see https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
# see https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/

# validating admission controller
$ cd validating-admission

$ make docker-build docker-push
$ make deploy-all
$ kubectl get all

$ kubectl deploy -f k8s/examples/nginx-ok.yaml
$ kubectl deploy -f k8s/examples/nginx-nok.yaml
$ kubectl delete -f k8s/examples/
$ make undeploy-all

# mutating admission controller
$ cd mutating-admission

$ make docker-build docker-push
$ make deploy-all
$ kubectl get all

$ kubectl deploy -f k8s/examples/nginx-ok.yaml
$ kubectl deploy -f k8s/examples/nginx-nok.yaml
$ kubectl delete -f k8s/examples/
$ make undeploy-all
```

## Maintainer

M.-Leander Reimer (@lreimer), <mario-leander.reimer@qaware.de>

## License

This software is provided under the MIT open source license, read the `LICENSE` file for details.

