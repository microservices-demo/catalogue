#!/usr/bin/env bash

set -ev

SCRIPT_DIR=$(dirname "$0")

if [[ -z "$GROUP" ]] ; then
    echo "Cannot find GROUP env var"
    exit 1
fi

if [[ -z "$COMMIT" ]] ; then
    echo "Cannot find COMMIT env var"
    exit 1
fi

if [[ "$(uname)" == "Darwin" ]]; then
    DOCKER_CMD=docker
else
    DOCKER_CMD="sudo docker"
fi
CODE_DIR=$(cd $SCRIPT_DIR/..; pwd)
echo $CODE_DIR

REPO=${GROUP}/$(basename catalogue);

$DOCKER_CMD build -t ${REPO}-dev $CODE_DIR/docker/catalogue;
$DOCKER_CMD create --name catalogue ${REPO}-dev;
$DOCKER_CMD cp catalogue:/app/main $CODE_DIR/docker/catalogue;
$DOCKER_CMD rm catalogue;
$DOCKER_CMD build -t ${REPO}:${COMMIT} -f $CODE_DIR/docker/catalogue/Dockerfile-release $CODE_DIR/docker/catalogue;
