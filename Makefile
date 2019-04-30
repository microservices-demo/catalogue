NAME = gcr.io/kavach-builds/sock-shop/catalogue
DBNAME = weaveworksdemos/catalogue-db

TAG=latest

INSTANCE = catalogue

.PHONY: default copy test

default: test

release:
	docker build -t $(NAME):$(TAG) -f ./docker/catalogue/Dockerfile .
	docker push $(NAME):$(TAG)

test: 
	GROUP=weaveworksdemos COMMIT=test ./scripts/build.sh
	./test/test.sh unit.py
	./test/test.sh container.py --tag $(TAG)

dockertravisbuild: build
	docker build -t $(NAME):$(TAG) -f docker/catalogue/Dockerfile-release docker/catalogue/
	docker build -t $(DBNAME):$(TAG) -f docker/catalogue-db/Dockerfile docker/catalogue-db/
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
	scripts/push.sh
