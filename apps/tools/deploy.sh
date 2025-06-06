#!/bin/bash

# constants
PORT=8086
IMAGE_NAME="tools"
DIRECTORY_NAME="tools"
CONTAINER_NAME="tools"

# prepare directory
cd ~/services/
rm -rf $DIRECTORY_NAME
git clone https://github.com/viveknathani/tools.git
cd $DIRECTORY_NAME

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t tools .
docker run -d --label tools=latest \
    --name=$CONTAINER_NAME \
    --network host \
    --pid=host \
    $IMAGE_NAME
