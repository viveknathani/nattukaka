FROM golang:1.18-bullseye
WORKDIR /usr/src/app
COPY . .
RUN apt-get update -y \
    && apt-get install -y \
    make \
    cron \
    && make build \
    && cp scripts/crontab /etc/cron.d/nattukaka \
    && chmod 0644 /etc/cron.d/nattukaka \
    && crontab /etc/cron.d/nattukaka
CMD ["make", "run"]
