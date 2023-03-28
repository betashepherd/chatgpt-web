MAIN_PKG := xray-client
DOCKER_HUB := ccr.ccs.tencentyun.com/fastapp
GIT_TAG := $(shell git describe --abbrev=0 --tags 2>/dev/null || echo 0.0.0)
GIT_TAG := v1.0.0
GIT_COMMIT_SEQ := $(shell git rev-parse --short HEAD 2>/dev/null || echo 000000)
GIT_COMMIT_CNT := $(shell git rev-list --all --count 2>/dev/null || echo 0)
VERSION := $(GIT_TAG).$(GIT_COMMIT_CNT).$(GIT_COMMIT_SEQ)
BUILD_TIME := $(shell TZ=UTC-8 date +"%Y%m%d%H%M%S")
FULL_VERSION := $(MAIN_PKG):$(GIT_TAG).$(GIT_COMMIT_CNT).$(GIT_COMMIT_SEQ)

mod:
	go mod tidy; go mod vendor

build:
	go build -tags=jsoniter -mod=vendor -ldflags "-s -w -X 'chatgpt-web/config.BuildVersion=$(FULL_VERSION)' -X 'chatgpt-web/config.BuildTime=$(BUILD_TIME)'" -o $(MAIN_PKG)

frontend:
	#cd frontend && npm install --registry=https://registry.npm.taobao.org && npm run build:prod && rm -rf ../dist && mv -f dist ../
	cd chat-new && cnpm install --registry=https://registry.npm.taobao.org && cnpm run build

docker-login:
	docker login ccr.ccs.tencentyun.com --username=100001261741 --password=ccr15652936798

docker:
	docker build . -t $(DOCKER_HUB)/$(FULL_VERSION)

docker-save: docker
	docker save $(DOCKER_HUB)/$(FULL_VERSION) | gzip > $(FULL_VERSION).tar.gz

docker-push:
	docker push $(DOCKER_HUB)/$(FULL_VERSION)

docker-clean:
	docker images | grep $(MAIN_PKG) | awk '{print $$3}' | xargs docker rmi -f
	echo "y" | docker image prune

version:
	echo $(DOCKER_HUB)/$(FULL_VERSION)

deploy:
	kubectl --kubeconfig .k3s.yaml config set-context --current --namespace $(MAIN_PKG)
	kubectl --kubeconfig .k3s.yaml set image deployment $(MAIN_PKG) $(MAIN_PKG)=ccr.ccs.tencentyun.com/fastapp/$(FULL_VERSION)

.PHONY: build frontend
