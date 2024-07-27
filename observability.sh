#!/bin/bash

docker run -d --name=grafana -p 3000:3000 grafana/grafana
docker run -d --name=loki -p 3100:3100 grafana/loki
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
