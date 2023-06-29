FROM ubuntu:20.04

COPY crawler /root/crawler
COPY config.json /root/config.json

WORKDIR /root

CMD ["./crawler", "-c", "./config.json", "-d", "./crawler.db"]
