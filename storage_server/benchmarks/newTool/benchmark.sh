#!/usr/local/bin/bash
# remember to: chmod +x benchmark.sh 


# Run go custom benchmark program 
CLIENTS=$1
REQUESTS=$2
MODE=$3
TYPE=$4


for i in $(seq 1 $CLIENTS)
    do
        nohup go run main.go $REQUESTS $MODE $i $TYPE >/dev/null 2>&1 &
        # go run main.go $REQUESTS $MODE $i
done