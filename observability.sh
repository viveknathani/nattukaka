#!/bin/bash

docker network create nattukaka-network
docker run -d --name=grafana -p 3000:3000 --network nattukaka-network grafana/grafana
docker run -d --name=loki -p 3100:3100 --network nattukaka-network grafana/loki
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
docker run -d --name=cadvisor -p 3200:3200 \
  --mount type=bind,source=/,target=/rootfs:ro \
  --mount type=bind,source=/var/run,destination=/var/run:ro \
  --network nattukaka-network \
  google/cadvisor:latest
docker run -d --name=prometheus -p 9090:9090 \
  -v /root/nattukaka/prometheus.yml:/etc/prometheus/prometheus.yml \
  --network nattukaka-network \
  prom/prometheus:latest