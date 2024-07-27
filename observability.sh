#!/bin/bash

docker network create nattukaka-network
docker run -d --name=grafana -p 3000:3000 --network nattukaka-network grafana/grafana
docker run -d --name=loki -p 3100:3100 --network nattukaka-network grafana/loki
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
