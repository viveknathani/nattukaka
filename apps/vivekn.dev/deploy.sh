#!/bin/bash

# constants
PORT=8081
IMAGE_NAME="vivekn.dev"
DIRECTORY_NAME="vivekn.dev"
CONTAINER_NAME="vivekn.dev"
ENV_FILE_PATH=~/environments/vivekn.dev

# check for env file
if [ ! -e "$ENV_FILE_PATH" ]; then
    echo ".env file does not exist for: " + $ENV_FILE_PATH
    exit 1
else
    echo ".env file found, proceeding!" 
fi

# prepare directory
mkdir -p services
cd ~/services/
rm -rf $DIRECTORY_NAME
git clone https://github.com/viveknathani/vivekn.dev.git
cd $DIRECTORY_NAME
cp $ENV_FILE_PATH .env

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t vivekn.dev .
docker run -d --label vivekn.dev=latest \
    --name=$CONTAINER_NAME \
    --network host \
    $IMAGE_NAME