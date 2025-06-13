#!/bin/bash

# constants
PORT=8087
IMAGE_NAME="chat"
DIRECTORY_NAME="chat"
CONTAINER_NAME="chat"
ENV_FILE_PATH=~/environments/chat

# prepare directory
cd ~/services/
rm -rf $DIRECTORY_NAME
git clone git@github.com:viveknathani/chat.git
cd $DIRECTORY_NAME
cp $ENV_FILE_PATH .env

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t chat .
docker run -d --label chat=latest \
    --name=$CONTAINER_NAME \
    --network host \
    $IMAGE_NAME
