#!/bin/bash

wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-7.x.list
apt update && apt install -y logstash filebeat

cp filebeat.yml /usr/share/filebeat/filebeat.yml
cp fishies.conf /usr/share/logstash/pipeline

systemctl restart logstash.service
systemctl restart filebeat.service
