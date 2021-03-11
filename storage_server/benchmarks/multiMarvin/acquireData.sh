#!/usr/local/bin/bash
# remember to: chmod +x acquireData.sh 

# gather data from FB server to LB server
read -p "Number of first bbchain machine: " FB
read -p "Number of last bbchain machine: " LB

# scp data to bbchain1
for i in $(seq $FB $LB)
    do
        ssh jmcad@bbchain$i.ux.uis.no -n 'cd Practice/storage_server/benchmarks/multiMarvin/&&for f in *.csv; do mv "$f" "${f%.csv}_$HOSTNAME.csv"; done; scp ./*.csv jmcad@bbchain1.ux.uis.no:Practice/storage_server/benchmarks/multiMarvin&&rm *.csv; exit'
done