FROM ubuntu:20.04

COPY crawler /root/crawler

WORKDIR /root

CMD ["./crawler", "-c", "./config.json", "-d", "./crawler.db"]
