FROM debian:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab

RUN apt update
RUN apt install -yq netcat iproute2
#RUN apt-get install iputils-ping

COPY ./src/ /src

# go compile
CMD ["src/main"]

# not needed anymore, replaced by start script
#ENTRYPOINT /bin/sh -c 'echo hello|nc -lvnp 80'
