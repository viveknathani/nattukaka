#!/bin/bash

# constants
PORT=8084
IMAGE_NAME="numero"
DIRECTORY_NAME="numero"
CONTAINER_NAME="numero"

# prepare directory
cd ~/services/
rm -rf $DIRECTORY_NAME
git clone https://github.com/viveknathani/numero.git
cd $DIRECTORY_NAME

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t numero .
docker run -d --label numero=latest \
    --name=$CONTAINER_NAME \
    --network host \
    --pid=host \
    $IMAGE_NAME
