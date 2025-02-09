#!/bin/bash

# constants
PORT=8083
IMAGE_NAME="workdiff"
DIRECTORY_NAME="workdiff"
CONTAINER_NAME="workdiff"
ENV_FILE_PATH=~/environments/workdiff

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
git clone https://github.com/viveknathani/workdiff.git
cd $DIRECTORY_NAME
cp $ENV_FILE_PATH .env

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t workdiff .
docker run -d --label workdiff=latest \
    --name=$CONTAINER_NAME \
    --network host \
    $IMAGE_NAME