docker-machine ssh dev "docker stop `$(docker ps -a -q)"
docker-machine ssh dev "docker rm `$(docker ps -a -q)"
