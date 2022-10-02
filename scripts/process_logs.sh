#!/bin/bash

cd ~/nattukaka && go run cmd/reporter/main.go > /var/log_stats.txt
