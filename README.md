# nattukaka

[![Go Report Card](https://goreportcard.com/badge/github.com/viveknathani/nattukaka)](https://goreportcard.com/report/github.com/viveknathani/nattukaka) 

nattukaka is a minmialistic web server that runs my stuff on the internet. It is designed in a highly personalized manner but you can easily fork this and run your own version. It currently runs on the cheapest VPS available on [Linode](https://www.linode.com). 

## features

- [x] serve my personal page ([vivekn.dev](https://vivekn.dev])) 
- [x] serve my blog ([vivekn.dev/blog](https://vivekn.dev/blog)) 
- [x] gated access to some parts ([vivekn.dev/login](https://vivekn.dev/login)) 
- [x] handle my todo list 
- [x] handle my notes 
- [x] send me the top 10 hacker news stories everyday at 7:00AM IST. 
- [ ] store and provide insights to my financial data 
- [ ] enable tagging on all posts (notes, blogs, etc.)

## requirements

- Go 1.18+ 
- PostgreSQL 
- Redis 
- [supercronic](https://github.com/aptible/supercronic) (runs my without hassle in docker containers)
- Docker (which will take care of setting up all of the above requirements for you, so just have this at the bare minimum) 
- nginx (optional but recommended, mostly managed by [certbot](https://certbot.eff.org/)). nginx would directly run on your host machine. 

## why this name? 

nattukaka is a famous character in an Indian TV show called Taarak Mehta Ka Ooltah Chashmah, who with his nephew, can get anything done for his employer. It sounds like a cool name for a do-it-all-for-me server! 

## why run your own server? 

- I love doing it! 
- I care about data privacy. Self-reliance feels good. 

## why not self-host existing open source software? 

Good option. I considered them but they're too bulky for my needs. nattukaka is lean. It does just as much as I want and is not bloated with a dozen things. Running multiple open source software is not trivial either. It is much lower headache to have almost everything run under one server. 

## setup info

- `.env` file requirements (this must be available on your VPS before running `docker compose`):
```bash
PORT=
JWT_SECRET=
TELEGRAM_API_KEY=
DATABASE_HOST=
DATABASE_PORT=
DATABASE_NAME=
DATABASE_USER=
DATABASE_PASSWORD=
REDIS_HOST=
REDIS_PORT=
```
- the `deploy.sh` script which is run on every push to master branch:
```bash
#!/bin/sh
echo "Removing previous copy"
sudo rm -rf nattukaka
echo "Getting latest"
git clone https://github.com/viveknathani/nattukaka.git
cp .env nattukaka/
cd nattukaka
docker compose build
docker compose up -d
```
## license

[MIT](./LICENSE)

