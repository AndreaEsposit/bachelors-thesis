#!/usr/local/bin/bash
# remember to: chmod +x launchBenchmark.sh; 


# run benchmark from FB server to LB server
read -p "Number of first bbchain machine: " FB
read -p "Number of last bbchain machine: " LB

read -p "Number of clients per machine:" CLIENTS
read -p "Number of requests per client:" REQUESTS
read -p "Mode:" MODE
read -p "Type:" TYPE

# Task 1: 
#   reset to head and pull master again to be sure that you have the latest build
pdsh -w andreaes@bbchain[$FB-$LB] "cd Practice/&&git reset --hard; git clean -f -d; git checkout HEAD; git pull"

# Task 2:
#   launhes the benchmark on all the bbchain machines
pdsh -w andreaes@bbchain[$FB-$LB] "cd Practice/storage_server/benchmarks/multi/&&chmod +x benchmark.sh&&./benchmark.sh $CLIENTS $REQUESTS $MODE $TYPE"

# ready the use of the acquire script
chmod +x acquireData.sh

# ready the use of the killAllProcesses script
chmod +x killAllProcesses.sh 

echo 'Wait a min or so. Before you run the acquireData script'