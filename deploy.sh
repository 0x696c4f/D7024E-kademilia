#!/bin/bash

# stop running containers 
docker kill $(docker ps -q -f ancestor=kadlab)
docker rm $(docker ps -a -q -f ancestor=kadlab)

# rebuild the container
docker build . -t kadlab

# start 50 containers
for i in {1..6};
do
	./start.sh
done
