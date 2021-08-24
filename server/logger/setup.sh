#!/bin/bash

apt update && apt install -y logstash filebeat

cp filebeat.yml /usr/share/filebeat/filebeat.yml
cp fishies.conf /usr/share/logstash/pipeline

systemctl restart logstash.service
systemctl restart filebeat.service
