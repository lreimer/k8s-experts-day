# Helm Operator Demo

```bash
$ take nginx-operator
$ operator-sdk init --plugins=helm.sdk.operatorframework.io/v1 \
                    --domain=cloud.qaware.de \
                    --project-name=nginx-operator \
                    --helm-chart=nginx-ingress \
                    --helm-chart-repo=https://helm.nginx.com/stable \
                    --version=v1alpha1
```
