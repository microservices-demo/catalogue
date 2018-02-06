#!/usr/bin/env sh

set -ev

export BUILD_VERSION="0.0.2-SNAPSHOT"
export BUILD_DATE=`date +%Y-%m-%dT%T%z`

SCRIPT_DIR=$(dirname "$0")

if [ -z "$GROUP" ] ; then
    echo "Cannot find GROUP env var"
    exit 1
fi

if [ -z "$COMMIT" ] ; then
    echo "Cannot find COMMIT env var"
    exit 1
fi

$(docker -v >/dev/null 2>&1)
if [ $? -eq 0 ]; then
    DOCKER_CMD=docker
else
    DOCKER_CMD=`sudo docker`
fi
CODE_DIR=$(cd $SCRIPT_DIR/..; pwd)
echo $CODE_DIR
 
mkdir -p $CODE_DIR/build
BUILD_DIR=$CODE_DIR/build

cp -r $CODE_DIR/docker $BUILD_DIR/
cp -r $CODE_DIR/images/ $BUILD_DIR/docker/catalogue/images/
cp -r $CODE_DIR/cmd/ $BUILD_DIR/docker/catalogue/cmd/
cp $CODE_DIR/*.go $BUILD_DIR/docker/catalogue/
mkdir -p $BUILD_DIR/docker/catalogue/vendor/ && \
cp $CODE_DIR/vendor/manifest $BUILD_DIR/docker/catalogue/vendor/

REPO=${GROUP}/$(basename catalogue);

$DOCKER_CMD build \
  --build-arg BUILD_VERSION=$BUILD_VERSION \
  --build-arg BUILD_DATE=$BUILD_DATE \
  --build-arg COMMIT=$COMMIT \
  -t ${REPO}:${COMMIT} \
  -f $BUILD_DIR/docker/catalogue/Dockerfile $BUILD_DIR/docker/catalogue;

$DOCKER_CMD build \
  -t ${REPO}-db:${COMMIT} \
  -f $BUILD_DIR/docker/catalogue-db/Dockerfile $BUILD_DIR/docker/catalogue-db/;

rm -rf $BUILD_DIR
