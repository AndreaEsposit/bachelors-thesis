#!/usr/local/bin/bash
# remember to: chmod +x cloneRepository.sh;
read -p "Your username:" USERNAME


pdsh -w $USERNAME@bbchain[2-30] "git clone git@github.com:AndreaEsposit/bachelors-thesis.git"