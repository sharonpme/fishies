#!/bin/bash

echo "Building and installing vernemq..."
cd vernemq
./install.sh
cd ..

echo "Installing service 'time'"
cd time
./setup.sh
cd ..

echo "Installing service 'scheduler'"
cd scheduler
./setup.sh
cd ..

echo "Installing service 'logger'"
cd logger 
./setup.sh
cd ..

echo "Setting up mDNS"
cd mdns
./setup.sh
cd ..

echo "Done!"
