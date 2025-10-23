#!/bin/bash

# constants
PORT=8086
IMAGE_NAME="tfcnc"
DIRECTORY_NAME="tradeforces"
CONTAINER_NAME="tfcnc"
ENV_FILE_PATH=~/environments/tfcnc

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
git clone git@github.com:viveknathani/tradeforces.git
cd $DIRECTORY_NAME
cp $ENV_FILE_PATH .env

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t tfcnc -f cnc.Dockerfile .
docker run -d --label tfcnc=latest \
    --name=$CONTAINER_NAME \
    --network host \
    $IMAGE_NAME
