export IMAGE="fr123k/jocker"

jocker: ## Start the jenkins in docker container short jocker.
	docker pull fr123k/jocker
	docker run -d -e BOOTSTRAP_SHARED_LIBRARY_GIT_REPO=fr123k/nmuxj -e BOOTSTRAP_SHARED_LIBRARY_GROOVY_FILE=jenkins/dsl/pipeline.groovy -e SEED_BRANCH=${TRAVIS_BRANCH} -p 8888:8080 --name jocker --rm ${IMAGE}

logs: ## Show the logs of the jocker container
	docker logs $(shell docker ps -f name=jocker -q) -f

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this help.
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
