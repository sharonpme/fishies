#!/bin/bash

SERVER_CONF="../server.conf"
USER_FILE="vmq.passwd"
INPUT_FILE="vernemq.conf.template"
OUTPUT_FILE="vernemq.conf"

usage() {
	echo "Usage: ./generate_config.sh [ -c <servo.conf> ] [ -u <user_file> ] [ -i <vernemq.conf.template> ] [ -o <vernemq.conf> ]
        -c Server configuration file. Default: ../server.conf
				-u User password file to be created. Default: vmq.passwd
				-i Input configuration template file. Default: vernemq.conf.template
				-o Output configuration file. Default: vernemq.conf"
	exit 2
}

while getopts 'c:ho:u:i:' c
do
	case $c in 
		c) SERVER_CONF=$OPTARG ;;
		h) usage ;;
		u) USER_FILE=$OPTARG ;;
		i) INPUT_FILE=$OPTARG ;;
		o) OUTPUT_FILE=$OPTARG ;;
	esac
done

mkdir -p $(dirname $OUTPUT_FILE)

get_var() {
	awk -F'=' "/^$2/ { print \$2 }" $1
}

HOST=$(get_var $SERVER_CONF MQTT_HOST)
PORT=$(get_var $SERVER_CONF MQTT_PORT)
USER=$(get_var $SERVER_CONF MQTT_USER)
PASS=$(get_var $SERVER_CONF MQTT_PASS)

# Generate password file
source /etc/profile
echo "$USER:$PASS" > $USER_FILE
vmq-passwd -U $USER_FILE

# Generate configuration file
tcp_listener_address="0.0.0.0:$PORT"
password_file=$USER_FILE 

eval "echo \"$(sed -r -e 's/`/\\`/g' $INPUT_FILE | sed -r -e 's/\$([^{])/\\$\1/g')\"" > $OUTPUT_FILE
