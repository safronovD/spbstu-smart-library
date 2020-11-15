#! /bin/bash
rm -r ./output 2> log.err 
docvar=local
if [[ ($1 == "remote") ]] 
then
    docvar=remote
fi
echo -e "******\n BUILDING APP \n******"
docker build --no-cache -t connector:builder ../ -f Dockerfiles/Dockerfile.$docvar 2> log.err
if [[ ($? -ne 0) ]]
then
{
    echo "Error : troubles building container"
    exit $ERRCODE
}
fi
echo -e "******\n DOWNLOADING DATA \n******"
docker container run --name dataLoader connector:builder ./connector.app 2> log.err
if [[ ($? -ne 0) ]]
then
{
    echo "Error : troubles running container"
    exit $ERRCODE
}
fi
docker container cp dataLoader:/go/connector/app/output/ ./output/ 2> log.err
if [[ ($? -ne 0) ]]
then
{
    echo "Error : troubles downloading data"
    exit $ERRCODE
}
fi
echo -e "******\n CLEANING \n******"
docker container rm -f dataLoader
docker images -q -f dangling=true | xargs -I ARGS docker rmi -qf ARGS

