# CNI Offload

# Set to valid docker registry
REPOSITORY := localhost:5000

fmt:
	go fmt ./...

vet:
	go vet ./...

clone-cilium:
	mkdir -p build;
	cd build && \
	rm -fr cilium && \
	git clone https://github.com/cilium/cilium.git && \
	cd cilium && git checkout c81fa2e447990c8834320b2c5401a4a759c76e55 && \
	git am ../../patches/cilium/*.patch;

clone-offload-cni:
	mkdir -p build;
	cd build && \
	rm -fr offload-cni && \
	git clone https://github.com/k8snetworkplumbingwg/sriov-cni offload-cni && \
	cd offload-cni && git checkout fca6591297c0e52b6522573bd367ccb7d6551fd0 && \
	git am ../../patches/offload-cni/*.patch;

cni-offload-agent:
	cd cmd/dpu/cniOffloadAgent; GOARCH=arm64 GOOS=linux go build ;

cilium-offload-cni:
	cd build/cilium/plugins/cilium-offload-cni; GOARCH=arm64 GOOS=linux make ;

dpu-components: cni-offload-agent cilium-offload-cni
	docker buildx build --network=host --push --tag $(REPOSITORY)/cni-offload-agent:latest -f Dockerfile.dpu --platform linux/arm64 .

offload-cni: clone-offload-cni
	cd build/offload-cni/ ; cd cmd/offload-cni; go build;

host-components: offload-cni
	docker build --network=host --push --tag $(REPOSITORY)/cni-offload-host:latest -f Dockerfile.host .

cilium:
	cd build/cilium/ ; ARCH=multi DOCKER_REGISTRY=$(REPOSITORY) make docker-cilium-image ;

build: fmt vet clone-cilium dpu-components host-components cilium

deploy-components:
	kubectl apply -f deployments/

deploy-cilium:
	kubectl apply -f build/cilium/install/kubernetes/cilium/cm.yaml
	sh deployments/install_cilium.sh $(REPOSITORY)

deploy-cilium-ipsec:
	sh deployments/install_cilium_ipsec.sh $(REPOSITORY)

deploy: deploy-components deploy-cilium

deploy-ipsec: deploy-components deploy-cilium-ipsec

undeploy:
	kubectl delete -f deployments/ --ignore-not-found=true
	kubectl delete -f build/cilium/install/kubernetes/cilium/cm.yaml --ignore-not-found=true
	helm delete cilium -n kube-system
	kubectl delete secret -n kube-system   cilium-ipsec-keys --ignore-not-found=true
