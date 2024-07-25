#!/bin/bash

# constants
PORT=8080
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
cp $ENV_FILE_PATH .env

# docker!
docker build -t teachyourselfmath .
docker run -dp $PORT:$PORT --label teachyourselfmath=latest --name $CONTAINER_NAME $IMAGE_NAME
