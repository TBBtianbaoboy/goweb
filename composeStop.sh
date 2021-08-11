#!/bin/bash

sudo docker-compose stop 
sudo docker container prune
sudo docker rmi go_web
