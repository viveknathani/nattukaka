# get the latest installation commands from https://caddyserver.com/

# Once done,
sudo mkdir -p /var/www/nattukaka
sudo nano /var/www/nattukaka/index.html
# add the following to begin with
# <!DOCTYPE html>
# <html lang="en">
# <head>
#     <meta charset="UTF-8">
#     <meta name="viewport" content="width=device-width, initial-scale=1.0">
#     <title>mmb</title>
#     <style>
#         body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
#     </style>
# </head>
# <body>
#     <p>hello</p>
# </body>
# </html>
sudo chown -R www-data:www-data /var/www/nattukaka
sudo chmod -R 755 /var/www/nattukaka

# now, edit the caddyfile, use your own.

# when done,
sudo systemctl restart caddy

# if using Cloudflare, make sure SSL mode is Full
