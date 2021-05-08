#!/usr/local/bin/bash
# remember to: chmod +x killAllProcesses.sh


read -p "User to kill:" USERNAME

pdsh -w $USERNAME@bbchain[1-30] "killall -u ${USERNAME}"
pdsh -w $USERNAME@bbchain[1-30] "ps -u ${USERNAME}" 