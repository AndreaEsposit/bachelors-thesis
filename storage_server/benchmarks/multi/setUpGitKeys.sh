#!/usr/local/bin/bash
# This script will copy your ssh key on all machines and then clone the Practice repository on all of them at the same time

# update build on every bbchain (git pull)

for i in $(seq 2 30);do 
    scp ~/.ssh/id_ed25519.pub andreaes@bbchain$i.ux.uis.no:~/.ssh
    scp ~/.ssh/id_ed25519 andreaes@bbchain$i.ux.uis.no:~/.ssh
done

pdsh -w andreaes@bbchain[2-30] "git clone git@github.com:AndreaEsposit/Practice.git"