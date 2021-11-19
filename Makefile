NAME = k8s-experts-day
VERSION = 1.0.0
GCP = gcloud
ZONE = europe-west1-b
K8S = kubectl
GITHUB_USER ?= lreimer

.PHONY: info

info:
	@echo "Kubernetes Experts Day"

prepare:
	@$(GCP) config set compute/zone $(ZONE)
	@$(GCP) config set container/use_client_certificate False

cluster:
	@echo "Create GKE Cluster"
	# --[no-]enable-basic-auth --[no-]issue-client-certificate

	@$(GCP) container clusters create $(NAME) --num-nodes=5 --enable-autoscaling --min-nodes=5 --max-nodes=10 --no-enable-autoupgrade
	@$(K8S) create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$$(gcloud config get-value core/account)
	@$(K8S) cluster-info

flux-bootstrap:
	@flux bootstrap github \
		--owner=$(GITHUB_USER) \
  		--repository=$(NAME) \
  		--branch=main \
  		--path=./flux2-demo/cluster \
		--components-extra=image-reflector-controller,image-automation-controller \
		--read-write-key \
  		--personal

gcloud-login:
	@$(GCP) auth application-default login

access-token:
	@$(GCP) config config-helper --format=json | jq .credential.access_token

destroy:
	@$(GCP) container clusters delete $(NAME) --async --quiet
