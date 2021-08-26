#!/bin/bash

# This script is written with a Debian system in mind since in all likelihood this is going to be run on a RPi.

apt update 
apt install -y curl gcc g++ git libsnappy-dev make libssl-dev perl sed tar wget

export BUILD_DIR=$HOME/tmp
mkdir -p $BUILD_DIR

# Building a very specific version of Erlang
cd $BUILD_DIR
wget https://github.com/erlang/otp/archive/OTP-23.3.4.4.tar.gz
tar -zxf OTP-23.3.4.4.tar.gz

export ERL_TOP=$BUILD_DIR/otp-OTP-23.3.4.4
cd $ERL_TOP
./otp_build autoconf
./configure --without-termcap
make
make install

# Retrieve and compile VerneMQ

cd $BUILD_DIR
git clone git://github.com/erlio/vernemq.git
cd vernemq
make rel

# Copy the VerneMQ directory somewhere nice and add it to the path
cp -r $BUILD_DIR/vernemq/_build/default/rel/vernemq /usr/bin
echo "" >> /etc/profile
echo "PATH=\"\$PATH:/usr/bin/vernemq/bin\"" >> /etc/profile
echo "export PATH" >> /etc/profile

source /etc/profile
rm -rf $BUILD_DIR
