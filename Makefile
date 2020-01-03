.PHONY: build shell jenkins local

time := $(shell date +'%Y%m%d-%H%M%S')
export IMAGE="nmuxj"

export DOCKER_COMMAND_LOCAL=docker run \
		-e JENKINS_SERVER="http://host.docker.internal:8888" \
		-e JENKINS_USER="admin" \
		-e APP_TOKEN="1234" \
		-e PORT=8080 \
		-p 8080:8080 \

build: ## Package the infra-hook go application into a go binary docker image.
	docker build -t $(IMAGE) -f Dockerfile .

local: build ## Build the provision hook service and start it
	$(DOCKER_COMMAND_LOCAL) -it -e JENKINS_TOKEN="$(token)" $(IMAGE)

travis: build ## Build the provision hook service and start it
	$(DOCKER_COMMAND_LOCAL) -d -e JENKINS_TOKEN="$(token)" $(IMAGE)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this help
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

jenkins:
	docker run -p 8888:8080 -p 50000:50000 -v $(PWD)/jenkins_home:/var/jenkins_home jenkins/jenkins:lts

test-hook:
	sleep 20
	curl -H "Content-Type:application/json" -H "Authorization:Bearer MTIzNA==" -d @test/jhook-param.json localhost:8080/provision/${time}/group/redis/total/2/vm/1
	curl -H "Content-Type:application/json" -H "Authorization:Bearer MTIzNA==" -d @test/jhook-param.json localhost:8080/provision/$(time)/group/redis/total/2/vm/2
	curl -H "Content-Type:application/json" -H "Authorization:Bearer MTIzNA==" -d @test/jhook-noparam.json localhost:8080/provision/$(time)/group/elasticsearch/total/1/vm/1
	sleep 20
	@curl -s http://admin:admin@localhost:8888/job/provision-redis/lastBuild/api/json | jq -r .result
	@curl -s http://admin:admin@localhost:8888/job/provision-elasticsearch/lastBuild/api/json | jq -r .result
