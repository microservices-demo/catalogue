NAME = weaveworksdemos/catalogue
DBNAME = weaveworksdemos/catalogue-db

TAG=$(TRAVIS_COMMIT)

INSTANCE = catalogue

.PHONY: default build copy

default: build

build:
	docker build -t $(NAME)-dev -f ./docker/catalogue/Dockerfile .

copy:
	docker create --name $(INSTANCE) $(NAME)-dev
	docker cp $(INSTANCE):/app/main $(shell pwd)/app
	docker rm $(INSTANCE)

release:
	docker build -t $(NAME) -f ./docker/catalogue/Dockerfile-release .

run:
	docker run --rm -p 8080:80 --name $(INSTANCE) $(NAME)

dockertravisbuild: build
	docker build -t $(NAME):$(TAG) -f docker/catalogue/Dockerfile-release docker/catalogue/
	docker build -t $(DBNAME):$(TAG) -f docker/catalogue-db/Dockerfile docker/catalogue-db/
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
	scripts/push.sh