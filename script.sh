#!/bin/sh

sudo docker rm -f containerdockerforum
sudo docker system prune -a

sudo docker build -t dockerforum .
sudo docker container run -p 8080:8080 -d --name containerdockerforum dockerforum
sudo docker image ls
sudo docker ps -a
sudo docker exec -it containerdockerforum /bin/bash
sudo ls -l