#!/bin/bash

docker network create nattukaka-network
docker run -d --name=grafana -p 3000:3000 grafana/grafana --network nattukaka-network
docker run -d --name=loki -p 3100:3100 grafana/loki --network nattukaka-network
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
