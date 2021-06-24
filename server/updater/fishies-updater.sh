#!/bin/sh

usage() {
	echo "Usage: ./fishies-updater.sh [ -c SERVER_CONF ] updated_firmware
         -c SERVER_CONF File with server connection details. Default: /etc/fishies/server.conf
				 -t TOPIC MQTT topic to publish to. Default: update"
	exit 2
}

SERVER_CONF=/etc/fishies/server.conf
TOPIC=update

while getopts 'c:ht:' c
do
	case $c in 
		c) SERVER_CONF=$OPTARG ;;
		h) usage ;;
		t) TOPIC=$OPTARG ;;
	esac
done

shift $((OPTIND - 1))

if [ -z $1 ]
then
	usage
fi

get_var() {
	awk -F'=' "/^$2/ { print \$2 }" $1
}

HOST=$(get_var ../server.conf MQTT_HOST)
PORT=$(get_var ../server.conf MQTT_PORT)
USER=$(get_var ../server.conf MQTT_USER)
PASS=$(get_var ../server.conf MQTT_PASS)

mosquitto_pub --topic $TOPIC -u $USER -P $PASS -h $HOST -p $PORT -m $(cat $1)
