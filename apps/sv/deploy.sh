#!/bin/bash

# constants
PORT=8085
IMAGE_NAME="sv"
DIRECTORY_NAME="sv"
CONTAINER_NAME="sv"
ENV_FILE_PATH=~/environments/sv

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
git clone git@github.com:viveknathani/sv.git
cd $DIRECTORY_NAME
cp $ENV_FILE_PATH .env

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t sv .
docker run -d --label sv=latest \
    --name=$CONTAINER_NAME \
    --network host \
    $IMAGE_NAME
