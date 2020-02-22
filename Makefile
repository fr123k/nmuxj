.PHONY: build shell jenkins local

time := $(shell date +'%Y%m%d-%H%M%S')
VERSION=1.0
export NAME=fr123k/nmuxj
export IMAGE="${NAME}:${VERSION}"
export LATEST="${NAME}:latest"
API_TOKEN=$(shell docker logs $(shell docker ps -f name=jocker -q) | grep 'Api-Token:' | tr ':' '\n' | tail -n +2)

export DOCKER_COMMAND_LOCAL=docker run \
		-e JENKINS_SERVER="http://host.docker.internal:8888" \
		-e JENKINS_USER="admin" \
		-e APP_TOKEN="1234" \
		-e PORT=8080 \
		-p 8080:8080 \

build: ## Package the infra-hook go application into a go binary docker image.
	docker build -t $(IMAGE) -f Dockerfile .

release: build ## Push docker image to docker hub
	docker tag ${IMAGE} ${LATEST}
	docker push ${NAME}

local: build ## Build the provision hook service and start it
	$(DOCKER_COMMAND_LOCAL) -it -e JENKINS_TOKEN="$(API_TOKEN)" --name nmuxj --rm $(IMAGE)

travis: build ## Build the provision hook service and start it
	$(DOCKER_COMMAND_LOCAL) -d -e JENKINS_SERVER="http://$(shell ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+'):8888" -e JENKINS_TOKEN="$(API_TOKEN)" --name nmuxj --rm $(IMAGE)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this help
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test-hook:
	docker logs $(shell docker ps -f name=jocker -q)
	curl -H "Content-Type:application/json" -H "Authorization:Bearer MTIzNA==" -d @test/jhook-param.json localhost:8080/provision/${time}/group/redis/total/2/vm/1
	curl -H "Content-Type:application/json" -H "Authorization:Bearer MTIzNA==" -d @test/jhook-param.json localhost:8080/provision/$(time)/group/redis/total/2/vm/2
	curl -H "Content-Type:application/json" -H "Authorization:Bearer MTIzNA==" -d @test/jhook-noparam.json localhost:8080/provision/$(time)/group/elasticsearch/total/1/vm/1
	docker logs $(shell docker ps -f name=nmuxj -q)
	sleep 20
	@curl -s -u admin:$(API_TOKEN) http://localhost:8888/job/provision-redis/lastBuild/api/json | jq -r .result | grep SUCCESS
	@curl -s -u admin:$(API_TOKEN) http://localhost:8888/job/provision-elasticsearch/lastBuild/api/json | jq -r .result | grep SUCCESS
