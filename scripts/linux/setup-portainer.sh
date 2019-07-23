#!/bin/sh
docker volume create portainer_data
docker run --restart=always --name portainer_srv -d -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data portainer/portainer
