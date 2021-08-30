#!/bin/bash

wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-7.x.list
apt update && apt install -y logstash filebeat

cp filebeat.yml /etc/filebeat/filebeat.yml
cp fishies.conf /etc/logstash/conf.d/fishies.conf

touch /etc/fishies/cfeed.log
chown logstash /etc/fishies/cfeed.log

touch /etc/fishies/ctime.log
chown logstash /etc/fishies/ctime.log

systemctl enable logstash.service
systemctl restart logstash.service
systemctl enable filenenat.service
systemctl restart filebeat.service
