# fishies server

The server system is centred around a distributed MQTT system ([VerneMQ](https://vernemq.com/)), which handles all the syncing between nodes. From there, various services interact with the queue to affect the state of the feeders.
