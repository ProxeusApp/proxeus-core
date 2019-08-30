FROM debian:stretch-slim

RUN apt-get update && apt-get install -y \
        ca-certificates \
    && apt-get clean && rm -rf /var/cache/* /var/lib/apt/lists/*

RUN update-ca-certificates --fresh

RUN mkdir -p /data/hosted

WORKDIR /app

COPY artifacts/server-docker /app/server

RUN chmod +x ./server

CMD ["./server"]
