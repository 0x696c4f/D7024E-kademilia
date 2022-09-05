#!/bin/sh

PORT=4000
COMMAND="./kademlia"
OPTIONS=""

if [ $(docker ps -q -f ancestor=kadlab|wc -l) -eq 0 ]
then
	# first container
	echo "Starting first container..."
	COMMAND="$COMMAND start"
else
	# join existing network
	# find ip
	IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps -q|shuf|head -n 1))
	echo "Join network via node $IP:$PORT..."
	COMMAND="$COMMAND join $IP $PORT"
fi
COMMAND="$COMMAND $OPTIONS"

# debug
COMMAND='echo hello | nc -lvnp 80'
docker run -d kadlab /bin/sh -c "$COMMAND"
