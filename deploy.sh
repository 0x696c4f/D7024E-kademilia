#!/bin/sh

# stop running containers 
docker kill $(docker ps -q -f ancestor=kadlab)
docker rm $(docker ps -a -q -f ancestor=kadlab)

# rebuild the container
docker build . -t kadlab

# start 50 containers
for ((i=0;i<50;i++));
do
	./start.sh
done
