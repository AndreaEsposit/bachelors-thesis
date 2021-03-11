#!/usr/local/bin/bash
# remember to: chmod +x cloneRepository.sh;

pdsh -w jmcad@bbchain[2-30] "git clone git@github.com:AndreaEsposit/Practice.git"