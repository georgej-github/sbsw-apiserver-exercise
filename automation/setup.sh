#!/bin/sh

DOCKER_IMAGE=test
DOCKER_REPO=test
IMAGE_VERSION=latest

#git pull https://github.com/georgej-github.com/sbsw-apiserver-exercise.git

#cd sbsw-apiserver-exercise

set -x

docker build -t $DOCKER_IMAGE -f ../docker/Dockerfile ../

docker tag $DOCKER_IMAGE $DOCKER_REPO/$DOCKER_IMAGE:$IMAGE_VERSION
docker push $DOCKER_REPO/$DOCKER_IMAGE:$IMAGE_VERSION

sed "s/%%IMAGEURL%%/${DOCKER_REPO}/${DOCKER_IMAGE}:${IMAGE_VERSION}/g" manifests/apiserver.yml
kubectl -f apply manifests/apiserver.yml



