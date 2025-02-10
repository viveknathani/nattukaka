#!/bin/bash

# packages
sudo apt update
sudo apt upgrade

# firewall and protection
sudo ufw allow 22
sudo ufw allow http
sudo ufw allow https
sudo ufw enable
sudo ufw status
sudo apt install fail2ban
sudo ufw allow from 192.168.1.0/24 to any port 8080:9000 proto tcp

# docker
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# postgresql
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql.service
sudo -i -u postgres
createuser --interactive # add viveknathani here
createdb viveknathani
exit
sudo adduser viveknathani
sudo -u viveknathani psql # verify that this is working by \conninfo
exit
# replace 16 with your version number (if needed)
sudo vi /etc/postgresql/16/main/postgresql.conf # change -> listen_addresses = '*'
sudo ufw allow 6379/tcp
sudo service postgresql restart

# redis
sudo apt update
sudo apt install redis-server
sudo vi /etc/redis/redis.conf # change -> 1) requirepass your_secure_password_here 2) bind 0.0.0.0 3) maxmemory-policy noeviction
sudo ufw allow 6379/tcp
sudo service redis-server restart

# clickhouse
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg
curl -fsSL 'https://packages.clickhouse.com/rpm/lts/repodata/repomd.xml.key' | sudo gpg --dearmor -o /usr/share/keyrings/clickhouse-keyring.gpg

ARCH=$(dpkg --print-architecture)
echo "deb [signed-by=/usr/share/keyrings/clickhouse-keyring.gpg arch=${ARCH}] https://packages.clickhouse.com/deb stable main" | sudo tee /etc/apt/sources.list.d/clickhouse.list
sudo apt-get update
sudo apt-get install -y clickhouse-server clickhouse-client
sudo ufw allow 9000/tcp
sudo ufw allow 8123/tcp
