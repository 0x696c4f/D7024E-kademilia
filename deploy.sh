docker build . -t kadlab
docker swarm init
docker stack deploy -c docker-compose.yml kademliaNodes
