FROM golang:latest

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
RUN apt install -yq apt-utils
RUN apt upgrade -y --fix-missing
RUN apt install -yq netcat iproute2 git
#RUN apt-get install iputils-ping

# go compile
ENV GOPATH=/kademlia/
ENV GO111MODULE=off
COPY . /kademlia
RUN go get github.com/gin-gonic/gin
WORKDIR /kademlia/src/D7024E
RUN go install

#ENTRYPOINT go run /src/main.go

# not needed anymore, replaced by start script
#ENTRYPOINT /bin/sh -c 'echo hello|nc -lvnp 80'
