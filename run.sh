#!/bin/bash

# stop postgresql service if it is running
sudo service postgresql stop

# Start the go_db service in the background
docker-compose up -d go_db

# Wait for a few seconds to make sure the go_db service is running
sleep 5

# Build and start the go-app service
docker-compose up --build go_app