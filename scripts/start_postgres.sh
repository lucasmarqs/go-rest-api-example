#!/usr/bin/env sh

volume_name=golang_todoapp

volume_exists=$(docker volume ls | grep $volume_name | wc -l)

if [ $volume_exists = '0' ]; then
  docker volume create $volume_name
fi

docker run --rm -d -p 5432:5432 \
  --mount source=$volume_name,target=/var/lib/postgresql/data \
  -e POSTGRES_USER=todoapp \
  -e POSTGRES_PASSWORD=todoapp \
  -e POSTGRES_DB=todoapp \
  --name todoapp-db \
  postgres:11.5
