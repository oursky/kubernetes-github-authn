REPO := nwik/kubernetes-gitlab-authn
VERSION := 0.4
IMAGE_NAME := $(REPO):$(VERSION)
GO_SRC_PATH := /go/src/github.com/$(REPO)
PORT := 8080

ifeq (1,${WITH_DOCKER})
DOCKER_RUN := docker run --rm -i \
	-v `pwd`:$(GO_SRC_PATH) \
	-w $(GO_SRC_PATH)
GO_RUN := $(DOCKER_RUN) golang:1.13-alpine
GLIDE_RUN := $(DOCKER_RUN) -e GLIDE_HOME=/root/.glide lwolf/golang-glide
endif

.PHONY: build
build:
	$(GO_RUN) CGO_ENABLED=0 go build -o _output/main main.go

.PHONY: vendor
vendor:
	$(GLIDE_RUN) glide install

.PHONY: clean
clean:
	rm -rf _output

.PHONY: docker-build
docker-build:
	#WITH_DOCKER=1 make build
	docker build --build-arg=https_proxy=${http_proxy} --build-arg=http_proxy=${http_proxy} -t $(IMAGE_NAME) .

docker-push:
	docker push $(IMAGE_NAME)

.PHONY: docker-run
docker-run:
	docker run -it --rm -p $(PORT):3000 $(IMAGE_NAME)
