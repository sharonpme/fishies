apt update && apt install -y avahi-daemon
cp mqtt.service /etc/avahi/services
hostname -I | head -n 1 | sed -e 's/$/mqtt/' >> /etc/avahi/hosts

systemctl enable avahi-daemon.service
systemctl start avahi-daemon.service
