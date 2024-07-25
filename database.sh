#!/bin/bash

# postgres
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
