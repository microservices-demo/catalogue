[![Build Status](https://travis-ci.org/microservices-demo/catalogue.svg?branch=master)](https://travis-ci.org/microservices-demo/catalogue) 
[![Coverage Status](https://coveralls.io/repos/github/microservices-demo/catalogue/badge.svg?branch=master)](https://coveralls.io/github/microservices-demo/catalogue?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/microservices-demo/catalogue)](https://goreportcard.com/report/github.com/microservices-demo/catalogue)
[![Microbadger](https://images.microbadger.com/badges/image/weaveworksdemos/catalogue.svg)](http://microbadger.com/images/weaveworksdemos/catalogue "Get your own image badge on microbadger.com")


# Catalogue
A microservices-demo service that provides catalogue/product information. 
This service is built, tested and released by travis.

## Bugs, Feature Requests and Contributing
We'd love to see community contributions. We like to keep it simple and use Github issues to track bugs and feature requests and pull requests to manage contributions.


### API Spec

Checkout the API Spec [here](http://microservices-demo.github.io/api/index?url=https://raw.githubusercontent.com/microservices-demo/catalogue/master/api-spec/catalogue.json)


### To build this service

#### Dependencies
```
go get -u github.com/FiloSottile/gvt
gvt restore
```

#### Go tools
In order to build the project locally you need to make sure that the repository directory is located in the correct
$GOPATH directory: $GOPATH/src/github.com/microservices-demo/catalogue/. Once that is in place you can build by running:

```
cd $GOPATH/src/github.com/microservices-demo/catalogue/cmd/cataloguesvc/
go build -o catalogue
```

The result is a binary named `catalogue`, in the current directory.

#### Docker
`docker-compose build`

### To run the service on port 8080

#### Go native

If you followed to Go build instructions, you should have a "catalogue" binary in $GOPATH/src/github.com/microservices-demo/catalogue/cmd/cataloguesvc/.
To run it use:
```
./catalogue
```

#### Docker
`docker-compose up`

### Run tests before submitting PRs
`make test`

### Check whether the service is alive
`curl http://localhost:8080/health`

### Use the service endpoints
`curl http://localhost:8080/catalogue`

### Push the service to Docker Container Registry
`GROUP=weaveworksdemos COMMIT=test ./scripts/push.sh`

## Test Zipkin

To test with Zipkin

```
docker-compose -f docker-compose-zipkin.yml build
docker-compose -f docker-compose-zipkin.yml up
```
It takes about 10 seconds to seed data

you should see it at:
[http://localhost:9411/](http://localhost:9411)

be sure to hit the "Find Traces" button.  You may need to reload the page.

when done you can run:
```
docker-compose -f docker-compose-zipkin.yml down
```
