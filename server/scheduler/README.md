# scheduler

A service to automatically push next feed timings to the MQTT queue on a schedule. For now, the schedules are defined and stored with cron strings, since it's a nice and easy format for storing schedule data. A simple web frontend is also provided to add and remove order schedules. The schedules are cached to disk with Badger. 
