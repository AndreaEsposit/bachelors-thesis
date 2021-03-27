#!/usr/local/bin/bash
# remember to: chmod +x killAllProcesses.sh

pdsh -w andreaes@bbchain[1-30] "killall -u andreaes"
pdsh -w andreaes@bbchain[1-30] "ps -u andreaes"