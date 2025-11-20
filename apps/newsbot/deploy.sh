#!/bin/bash

# constants
PORT=8086
IMAGE_NAME="newsbot"
DIRECTORY_NAME="newsbot"
CONTAINER_NAME="newsbot"
ENV_FILE_PATH=~/environments/newsbot

# check for env file
if [ ! -e "$ENV_FILE_PATH" ]; then
    echo "Env file does not exist for: " + $ENV_FILE_PATH
    exit 1
else
    echo ".env file found, proceeding!" 
fi

# prepare directory
cd ~/services/
rm -rf $DIRECTORY_NAME
git clone https://github.com/viveknathani/newsbot.git
cd $DIRECTORY_NAME
cp $ENV_FILE_PATH .env

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t newsbot .
docker run -d --label newsbot=latest \
    --name=$CONTAINER_NAME \
    --network host \
    $IMAGE_NAME
