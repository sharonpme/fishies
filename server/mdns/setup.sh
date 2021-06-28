apt update && apt install -y avahi-daemon
cp mqtt.service /etc/avahi/services

systemctl enable avahi-daemon.service
systemctl start avahi-daemon.service
