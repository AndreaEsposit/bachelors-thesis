#!/usr/bin/bash

# Run go custom benchmark program 
read -p "Number of clients:" CLIENTS
read -p "Number of requests:" REQUESTS
read -p "Mode:" MODE

for i in $(seq 1 $CLIENTS)
    do
        nohup go run main.go $REQUESTS $MODE $i >/dev/null 2>&1 &
done

# Acquire data
# sleep 1m
# python3 analyse.py