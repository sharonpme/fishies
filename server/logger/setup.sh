#!/bin/bash

wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-7.x.list
apt update && apt install -y logstash filebeat

cp filebeat.yml /etc/filebeat/filebeat.yml
cp fishies.conf /etc/logstash/conf.d/fishies.conf

touch /etc/fishies/fishies.log
chown logstash /etc/fishies/fishies.log

systemctl restart logstash.service
systemctl restart filebeat.service
