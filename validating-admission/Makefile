# Image URL to use all building/pushing image targets
IMG ?= lreimer/validating-admission:latest

update-certs:
	./gen_certs.sh
	kubectl delete secret validating-admission-cert --ignore-not-found
	kubectl create secret generic validating-admission-cert -n default --from-file=key.pem=certs/validating-key.pem --from-file=cert.pem=certs/validating-crt.pem

deploy:
	kubectl apply -f manifest.yaml

undeploy:
	kubectl delete -f manifest.yaml --ignore-not-found

# Build the docker image
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push: docker-build
	docker push ${IMG}
