WEBHOOK_SERVICE?=validating-admission-service
NAMESPACE?=default
CONTAINER_REPO?=lreimer/validating-admission
CONTAINER_VERSION?=1.0.3
CONTAINER_IMAGE=$(CONTAINER_REPO):$(CONTAINER_VERSION)

.PHONY: deploy-all deploy-csr deploy-webhook undeploy-all docker-build docker-push

deploy-all: deploy-csr deploy-webhook

deploy-csr:
	kubectl apply -k k8s/csr/
	@echo Waiting for cert creation ...
	@sleep 15
	kubectl certificate approve $(WEBHOOK_SERVICE).$(NAMESPACE)

deploy-webhook:
	(cd k8s/deployment && kustomize edit set image CONTAINER_IMAGE=$(CONTAINER_IMAGE))
	kubectl apply -k k8s/deployment/

undeploy-all:
	kubectl delete -k k8s/deployment/ --ignore-not-found
	kubectl delete -k k8s/csr/ --ignore-not-found

# Build the docker image
docker-build:
	docker build . -t ${CONTAINER_IMAGE}

# Push the docker image
docker-push: docker-build
	docker push ${CONTAINER_IMAGE}
