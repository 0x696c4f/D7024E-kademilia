#!/bin/sh

PORT=4000
COMMAND=""
OPTIONS=""

if [ $(docker ps -q -f ancestor=kadlab|wc -l) -eq 0 ]
then
	# first container
	COMMAND="$COMMAND start"
else
	# join existing network
	# find ip
	IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps -q|shuf|head -n 1))
	COMMAND="$COMMAND join $IP $PORT"
fi
COMMAND="$COMMAND $OPTIONS"
docker run kadlab $COMMAND
