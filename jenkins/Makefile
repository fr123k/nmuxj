export IMAGE="denkins"

build: ## Build the jenkins in docker image.
	docker build -t $(IMAGE) -f Dockerfile .

denkins: build ## Start the jenkins in docker container short denkins. 
	docker run -d -p 8888:8080 ${IMAGE}

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this help.
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
