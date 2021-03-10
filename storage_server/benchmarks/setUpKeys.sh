#!/bin/bash

# check if is running as root
uid=$(id -u)
[ $uid -ne 0 ] && { echo "This script should be executed as root."; exit 1; }

user=$1
pubkey_file=$2
stud_home="/home/stud/$user"
authorized_keys="${stud_home}/.ssh/authorized_keys"

[[ "$#" -ne 2 ]] && { echo "Usage : setup-sshpwless.sh user pubkey_file_path"; exit 1; }

# update authorized keys on bbchain1
cat ${pubkey_file} >> ${authorized_keys}

# copy to other nodes
for i in {2..30} ; do
    echo "#bbchain$i"
    scp ${authorized_keys} bbchain$i:${stud_home}/.ssh/authorized_keys
    ssh bbchain$i chmod 640 ${authorized_keys}
done