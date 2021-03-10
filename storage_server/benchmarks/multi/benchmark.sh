#!/usr/local/bin/bash

# if it doesn't work: chmod +x benchmark.sh

# update build on every bbchain (git pull)
firstBbchain = 20
lastBbchain = 30



for i in $(seq $firstBbchain $lastBbchain);do 
    ssh jmcad@bbchain{$i}.ux.uis.no
    cd Practice/
    git reset --hard HEAD
    git pull


# # Run go custom benchmark program 
# read -p "Number of clients:" CLIENTS
# read -p "Number of requests:" REQUESTS
# read -p "Mode:" MODE

# for i in $(seq 1 $CLIENTS)
#     do
#         nohup go run main.go $REQUESTS $MODE $i >/dev/null 2>&1 &
# done

# # Acquire data
# # sleep 1m
# # python3 analyse.py