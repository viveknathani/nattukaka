#!/bin/bash

cd ~/nattukaka && go run cmd/processor/main.go > /var/log_stats.txt
rm /var/logs.txt
