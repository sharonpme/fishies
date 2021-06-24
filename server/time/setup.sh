apt update && apt install -y mosquitto-clients date

cp fishies-timer.service /etc/systemd/system/fishies-timer.service
cp fishies-timer-timer.service /etc/systemd/system/fishies-timer-timer.service

mkdir -p /etc/fishies
cp ../server.conf /etc/fishies/server.conf

systemctl daemon-reload
systemctl restart fishies-timer-timer.service
