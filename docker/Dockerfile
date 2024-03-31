FROM ubuntu:20.04

RUN apt update && apt install ca-certificates -y
COPY crawler /root/crawler

WORKDIR /root

CMD ["./crawler", "-c", "./config.json", "-d", "./crawler.db"]
