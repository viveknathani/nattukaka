sudo apt update
sudo apt install nginx

sudo cp ~/nattukaka/services/nattukaka.dev/nginx.conf /etc/nginx/sites-available/nattukaka.dev
sudo ln -s /etc/nginx/sites-available/nattukaka.dev /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
