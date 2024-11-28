.PHONY: build test run stop

# Build the docker image
build:
	docker build -t my-indexer .

# Run tests in docker container
test: build
	docker run --rm my-indexer

# Run the application
run: build
	docker run -d --name my-indexer-app my-indexer

# Stop the running container
stop:
	docker stop my-indexer-app || true
	docker rm my-indexer-app || true
