#!/usr/local/bin/bash
# if it doesn't work: chmod +x benchmark.sh 

FB=20
LB=30

# Task 1: 
#   reset to head and pull master again to be sure that you have the latest build
pdsh -w andreaes@bbchain[$FB-$LB] "cd Practice/&&git reset --hard HEAD; git pull"



echo HELLO!

# Run go custom benchmark program 
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