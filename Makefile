NAME = weaveworksdemos/catalogue
DBNAME = weaveworksdemos/catalogue-db

TAG=$(TRAVIS_COMMIT)

INSTANCE = catalogue

.PHONY: default copy test

default: build

pre:
	go get -u github.com/FiloSottile/gvt

deps: pre
	gvt restore

rm-deps:
	rm -rf vendor

copy:
	docker create --name $(INSTANCE) $(NAME)-dev
	docker cp $(INSTANCE):/app/main $(shell pwd)/app
	docker rm $(INSTANCE)

release:
	docker build -t $(NAME) -f ./docker/catalogue/Dockerfile-release .

test: 
	GROUP=weaveworksdemos COMMIT=test ./scripts/build.sh
	./test/test.sh unit.py
	./test/test.sh container.py --tag $(TAG)

clean: cleandocker
	# rm -rf bin
	rm -rf docker/user/bin
	rm -rf vendor

dockertravisbuild: build
	docker build -t $(NAME):$(TAG) -f docker/catalogue/Dockerfile-release docker/catalogue/
	docker build -t $(DBNAME):$(TAG) -f docker/catalogue-db/Dockerfile docker/catalogue-db/
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
	scripts/push.sh

build: deps
	mkdir -p bin
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/$(INSTANCE) cmd/cataloguesvc/main.go
