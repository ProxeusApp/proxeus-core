#!/bin/bash

trap 'kill 0' INT

export RS=${RS:-rs0}
export MONGO_PORT=${MONGO_PORT:-27017}
export MONGO0_PORT=${MONGO0_PORT:-${MONGO_PORT}}
export MONGO1_PORT=${MONGO1_PORT:-$((${MONGO_PORT} + 1))}
export MONGO2_PORT=${MONGO2_PORT:-$((${MONGO_PORT} + 2))}

mkdir -p /data/db/mongo{0,1,2}

echo "Starting instances on ports ${MONGO0_PORT} ${MONGO1_PORT} ${MONGO2_PORT}..."
docker-entrypoint.sh --dbpath /data/db/mongo0 --bind_ip_all --replSet ${RS} --port ${MONGO0_PORT} & 
docker-entrypoint.sh --dbpath /data/db/mongo1 --bind_ip_all --replSet ${RS} --port ${MONGO1_PORT} &
docker-entrypoint.sh --dbpath /data/db/mongo2 --bind_ip_all --replSet ${RS} --port ${MONGO2_PORT} &


echo "Waiting for startup..."
until mongo --host localhost:${MONGO0_PORT} --eval 'quit(db.runCommand({ ping: 1 }).ok ? 0 : 2)' &>/dev/null; do
  printf '.'
  sleep 1
done

echo "Started."

echo "Configuring replica set"
mongo --host localhost:${MONGO0_PORT} <<EOF
config={"_id":"${RS}","members":[{"_id":0,"host":"localhost:${MONGO0_PORT}"},{"_id":1,"host":"localhost:${MONGO1_PORT}"},{"_id":2,"host":"localhost:${MONGO2_PORT}"}]}
rs.initiate(config)
EOF


echo "Replica set status:"
mongo --host localhost:${MONGO0_PORT} <<EOF
rs.status()
EOF

wait


