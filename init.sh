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
