#!/bin/sh

HOME_DIR=`pwd`
UTILS_DIR="$HOME_DIR/lib/utils"

echo "Waiting for ES to be running"
cd $UTILS_DIR
./wait_for_es.sh

echo "Downloading data"
cd $HOME_DIR
./connector.app

echo "Creating pattern and graphs"
cd $UTILS_DIR
python3 ./graphs_import.py