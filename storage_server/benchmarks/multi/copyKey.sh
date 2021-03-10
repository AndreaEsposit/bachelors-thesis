#!/usr/local/bin/bash

# update build on every bbchain (git pull)
FB=2
LB=30



for i in $(seq $FB $LB);do 
    scp ~/.ssh/id_ed25519.pub andreaes@bbchain{$i}.ux.uis.no:~/.ssh
    scp ~/.ssh/id_ed25519 andreaes@bbchain{$i}.ux.uis.no:~/.ssh
done