#!/usr/local/bin/bash
# remember to: chmod +x acquireData.sh 

# gather data from FB server to LB server
read -p "Number of first bbchain machine: " FB
read -p "Number of last bbchain machine: " LB

# scp data to bbchain1
for i in $(seq $FB $LB)
    do
        ssh andreaes@bbchain$i.ux.uis.no -n 'cd Practice/storage_server/benchmarks/multi/&&for f in *.csv; do mv "$f" "${f%.csv}_$HOSTNAME.csv"; done; scp ./*.csv andreaes@bbchain1.ux.uis.no:Practice/storage_server/benchmarks/multi&&rm *.csv; exit'
done

# Runs the python code to analyze the data
python3 analyse.py

# cat the file
cat result.json

# remove the file
rm result.json