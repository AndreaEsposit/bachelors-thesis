#!/usr/bin/bash

# update build on every bbchain (git pull)
firstBbchain = 2
lastBbchain = 30



for i in $(seq $firstBbchain $lastBbchain);do 
    scp ~/.ssh/id_ed25519.pub andreaes@bbchain{$i}.ux.uis.no:~/.ssh
    scp ~/.ssh/id_ed25519 andreaes@bbchain{$i}.ux.uis.no:~/.ssh
