apt update && apt install -y golang-go

go build

mkdir -p /etc/fishies
cp scheduler /etc/fishies/scheduler
cp ../server.conf /etc/fishies/server.conf
cp scheduler.conf /etc/fishies/scheduler.conf
cp -r pages /etc/fishies/page

cp scheduler.service /etc/systemd/system/fishies-scheduler.service

systemctl daemon-reload
systemctl restart fishies-scheduler.service
