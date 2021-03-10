#!/usr/local/bin/bash
# remember to: chmod +x benchmark.sh 


# Run go custom benchmark program 
SERVERNUMBER = $1
CLIENTS = $2
REQUESTS = $3
MODE = $4


for i in $(seq 1 $CLIENTS)
    do
        nohup go run main.go $REQUESTS $MODE $SERVERNUMBER $i >/dev/null 2>&1 &
done
