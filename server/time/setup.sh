apt update && apt install -y mosquitto-clients

cp fishies-timer.service /etc/systemd/system/fishies-timer.service
cp fishies-timer.timer /etc/systemd/system/fishies-timer.timer

mkdir -p /etc/fishies
cp ../server.conf /etc/fishies/server.conf

systemctl daemon-reload
systemctl restart fishies-timer.timer
