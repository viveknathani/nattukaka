#!/bin/bash

sudo apt update
sudo apt install nginx

sudo mkdir -p /var/www/nattukaka.dev
sudo cp /root/nattukaka/services/nattukaka.dev/index.html /var/www/nattukaka.dev/
sudo cp /root/nattukaka/services/nattukaka.dev/dashboard.png /var/www/nattukaka.dev/
sudo chmod -R 755 /var/www/nattukaka.dev
sudo chown -R www-data:www-data /var/www/nattukaka.dev
sudo cp ~/nattukaka/services/nattukaka.dev/nginx.conf /etc/nginx/sites-available/nattukaka.dev
sudo ln -s /etc/nginx/sites-available/nattukaka.dev /etc/nginx/sites-enabled/

sudo cp ~/nattukaka/services/teachyourselfmath/nginx.conf /etc/nginx/sites-available/teachyourselfmath
sudo ln -s /etc/nginx/sites-available/teachyourselfmath /etc/nginx/sites-enabled/

sudo cp ~/nattukaka/services/vivekn.dev/nginx.conf /etc/nginx/sites-available/vivekn.dev
sudo ln -s /etc/nginx/sites-available/vivekn.dev /etc/nginx/sites-enabled/

sudo cp ~/nattukaka/services/workdiff/nginx.conf /etc/nginx/sites-available/workdiff
sudo ln -s /etc/nginx/sites-available/workdiff /etc/nginx/sites-enabled/

sudo cp ~/nattukaka/services/grafana/nginx.conf /etc/nginx/sites-available/grafana
sudo ln -s /etc/nginx/sites-available/grafana /etc/nginx/sites-enabled/

sudo nginx -t
sudo systemctl restart nginx
