#!/bin/bash 

./setup.sh


./generate_config.sh -c ../server.conf -o /usr/bin/vernemq/etc/vernemq.conf -u /usr/bin/vernemq/etc/vmq.passwd
cp vmq.acl /usr/bin/vernemq/etc/vmq.acl
cp ./run_vernemq.sh /usr/bin/vernemq/bin
