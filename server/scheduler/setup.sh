apt update && apt install -y golang-go

go build

mkdir -p /etc/fishies
cp scheduler /etc/fishies

