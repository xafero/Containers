#!/bin/sh
docker rmi --force $(docker images -f "dangling=true" -q)
docker container prune --force
docker system prune --force
docker volume prune --force
