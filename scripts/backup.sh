#!/bin/bash
cd /var
sudo tar cvzf nattukaka/memory`date +"%Y%m%d"`.tar.gz memory.txt
sudo tar cvzf nattukaka/health`date +"%Y%m%d"`.tar.gz health.txt
sudo rm /var/memory.txt
sudo rm /var/health.txt
