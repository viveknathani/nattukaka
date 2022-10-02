#!/bin/bash

echo $(date "+%s"),$(ps -C nattukaka -o rss --no-headers) >> /var/memory.txt;
echo "$(date),$(curl -s -o /dev/null -I -w "%{http_code}\n" https://vivekn.dev/health)" >> /var/health.txt