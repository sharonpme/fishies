#!/bin/bash 

./setup.sh

./generate_config.sh -c ../server.conf -o /usr/bin/vernemq/etc/vernemq.conf -u /usr/bin/vernemq/etc/vmq.passwd
cp vmq.acl /usr/bin/vernemq/etc/vmq.acl

cp vernemq.service /etc/systemd/system/vernemq.service
systemctl daemon-reload
systemctl enable vernemq.service
systemctl start vernemq.service
