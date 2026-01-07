#!/bin/bash

# constants
PORT=8082
IMAGE_NAME="teachyourselfmath"
DIRECTORY_NAME="teachyourselfmath"
CONTAINER_NAME="teachyourselfmath"
ENV_FILE_PATH=~/environments/teachyourselfmath

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
git clone https://github.com/viveknathani/teachyourselfmath.git
cd $DIRECTORY_NAME
cp $ENV_FILE_PATH .env

# docker!
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
docker build -t teachyourselfmath .
docker run -d --label teachyourselfmath=latest \
    --name=$CONTAINER_NAME \
    --network host \
    $IMAGE_NAME
