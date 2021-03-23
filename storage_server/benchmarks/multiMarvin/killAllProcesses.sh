#!/usr/local/bin/bash
# remember to: chmod +x killAllProcesses.sh

pdsh -w jmcad@bbchain[1-30] "killall -u jmcad"

pdsh -w jmcad@bbchain[1-30] "ps -u jmcad"