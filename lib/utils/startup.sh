#!/bin/sh

HOME_DIR=`pwd`
UTILS_DIR="$HOME_DIR/utils"

echo "Waiting for ES to be ready"
cd $UTILS_DIR
./waiting_for_es.sh

echo "Downloading data"
cd $HOME_DIR
./connector.app

echo "Creating pattern and graphs"
cd $UTILS_DIR
python3 ./graphs_import.py

echo "Data loading have compleated. You can close the proccess. (Ctrl+C)"