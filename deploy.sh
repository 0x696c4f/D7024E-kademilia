docker build . -t kadlab
docker swarm leave --force
docker swarm init
docker stack deploy -c docker-compose.yml kademliaNodes
