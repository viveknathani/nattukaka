#!/bin/bash

docker network create nattukaka-network
docker run -d --name=grafana -p 3000:3000 --network nattukaka-network grafana/grafana
docker run -d --name=loki -p 3100:3100 --network nattukaka-network grafana/loki
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions

docker run -d \
  --name=node_exporter \
  --network=nattukaka-network \
  --volume="/proc:/host/proc:ro" \
  --volume="/sys:/host/sys:ro" \
  --volume="/:/rootfs:ro" \
  quay.io/prometheus/node-exporter \
  --path.procfs=/host/proc \
  --path.sysfs=/host/sys \
  --collector.filesystem.ignored-mount-points "^/(sys|proc|dev|host|etc)($|/)"

docker run -d --name=cadvisor -p 8080:8080 \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:rw \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --network nattukaka-network \
  gcr.io/cadvisor/cadvisor:latest

mkdir -p /srv/prometheus/data
docker run -d --name=prometheus \
  -p 9090:9090 \
  -v /root/nattukaka/prometheus.yml:/etc/prometheus/prometheus.yml \
  -v /srv/prometheus/data:/prometheus \
  --network nattukaka-network \
  prom/prometheus
