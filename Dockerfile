FROM debian:stable

RUN apt-get update && apt-get install -y \
        ca-certificates \
    && apt-get clean && rm -rf /var/cache/* /var/lib/apt/lists/*

RUN update-ca-certificates --fresh

RUN mkdir -p /data/hosted

WORKDIR /app

COPY artifacts/proxeus-docker /app/proxeus

COPY demo/ /app/demo/

RUN chmod +x ./demo/restore-demo.sh

RUN chmod +x ./proxeus

CMD ["./proxeus"]
