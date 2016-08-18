NAME = weaveworksdemos/catalogue
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
