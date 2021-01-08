#!/bin/sh

#Using for waiting for es availability
#TODO: Remake
echo -ne "\033[33m\rWaiting --- "
wget -q http://elasticsearch:9200

while [ $? != 0 ]
do
    sleep 10
    echo -ne "\033[33m\rWaiting --- "
    wget -q http://elasticsearch:9200
done

echo -e "\033[32mDone --- Connection successed"
rm ./index.html
