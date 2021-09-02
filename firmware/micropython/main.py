from umqtt.simple import MQTTClient
import os
import machine
import time
import network
import ubinascii
import uasyncio
from mdns_client import Client
from mdns_client.service_discovery.txt_discovery import TXTServiceDiscovery

CLIENT_ID = 'F' + ubinascii.hexlify(machine.unique_id()).decode("utf-8") # prefix with a letter so logstash doesn't die on leading zeroes

WIFI_NAME = 'Pi3-AP'
WIFI_PASSWORD = 'raspberry'

# pins
conditioning_led_pin = 14 # D5
servo_pin = 2 # D4

servo_frequency = 50 # Hertz
servo = machine.PWM(machine.Pin(servo_pin), freq=servo_frequency)
food_hole_offset = 30

# servo control
servo_min_pos = -90
servo_max_pos = 90
servo_min_pos_ms = 1
servo_max_pos_ms = 2

duty = 1 / servo_frequency * 1000 # maximum duty cycle (milliseconds)
duty_max = 1023
duty_min = 0

def servo_pos(servo, degrees):
    degrees = (degrees - servo_min_pos) * ((servo_max_pos_ms - servo_min_pos_ms) - (servo_max_pos - servo_min_pos)) + servo_min_pos_ms # Converted to duty cycle (ms)
    degrees = degrees * ((duty_max - duty_min) / duty) + duty_min # Converted it to that weird arbitary analog value thing
    servo.duty(degrees)

# init servo (just in case)
servo_pos(servo, 0)

# constants
conditioning_time = 3 # seconds
## Default time to wait if there are no new orders
default_waiting_time = 3600 # 1 hour = 60 * 60 seconds
## How long will the unit deep sleep after feeding?
postfeedsleepms = 60000 # milliseconds

# topics
topics_prefix = "orders/"
master_topic =  topics_prefix + "all"
sub_topic = topics_prefix + "feeder"

ignore_master = False
try:
    f = open("ignore-master", "r")
    ignore_master = True
    f.close()
except:
    ignore_master = False

## conditioning LED = G14
lcond = machine.Pin(conditioning_led_pin, machine.Pin.OUT)

# configure RTC.ALARM0 to be able to wake the device
rtc = machine.RTC()
rtc.irq(trigger=rtc.ALARM0, wake=machine.DEEPSLEEP)

# mqtt callback
def sub_cb(topic, msg):
    topic = topic.decode("utf-8")
    print(topic)

    if topic == "time":
        print("time topic")
        global ttime
        ttime = int(msg.decode("utf-8"))
        print(ttime)

        strtosend = CLIENT_ID + "," + str(ttime)
        c.publish(b"ctime", strtosend)
    elif topic == master_topic:
        global orders
        if not ignore_master:
            orders = int(msg.decode("utf-8"))
            print(orders)
    elif topic == sub_topic:
        global orders
        orders = int(msg.decode("utf-8"))
        print(orders)

        f = open("ignore-master", "w")
        f.close()
    elif topic == "status":
        if msg.decode("utf-8") == "updatetime": # time for an update
            time.sleep(1)
            c.connect(clean_session=True)
            c.subscribe(b"update")
            time.sleep(1)

            print("checking for message")
            gc.collect()
            c.check_msg()
    elif topic == "update":
        print(msg)
        gc.collect()
        with open('main.py', 'wb') as fd:
            fd.write(msg)
        print("Updated main.py, resetting")
        machine.reset()
        time.sleep(3)

# setup the network:
wlan = network.WLAN(network.STA_IF)
ap_if = network.WLAN(network.AP_IF)
if ap_if.active():
    ap_if.active(False)

wlan.active(True)

print('connecting to network...')
wlan.connect('Pi3-AP', 'raspberry')
while not wlan.isconnected():
    machine.idle()

print(wlan.ifconfig())

# get MQTT address via mDNS
SERVER = "192.168.0.13"  ## the Pi3-AP
loop = uasyncio.get_event_loop()
client = Client(wlan.ifconfig()[0])
discovery = TXTServiceDiscovery(client)

global discoveries
async def discover_once():
    discoveries = await discovery.query_once("_mqtt", "_tcp")

loop.run_until_complete(discover_once())

if len(discoveries.ips) == 0:
    print('no MQTT service found via mDNS, deepsleeping for 5 minutes')
    rtc.alarm(rtc.ALARM0, 300000) # 5 * 60 * 1000ms

global mqtt_ip
for ip in discoveries.ips:
    mqtt_ip = ip
    break

# setup MQTT client
c = MQTTClient(CLIENT_ID, mqtt_ip)
c.set_callback(sub_cb)

# Check for an update
print(machine.reset_cause())
if (machine.reset_cause() == machine.HARD_RESET) | (machine.reset_cause() == 5):
    print("Hard reset detected, checking for update")
    c.connect(clean_session=True)
    c.subscribe(b"status")
    time.sleep(1)
    print("checking for message")
    c.check_msg()

# Setup time, global time variables, and check time
DEFAULT_TIME = 90000
global ttime
ttime = DEFAULT_TIME

# Get the time
timeouti = 0
while ttime == DEFAULT_TIME:
    try:
        c.connect(clean_session=True)
        c.subscribe(b"time")
        while ttime == DEFAULT_TIME:
            timeouti += 1
            c.check_msg()
            time.sleep(1)
            print(timeouti)
            if timeouti == 20:
                print("can't connect to server, deep sleeping for 5 min and will try again")
                rtc.alarm(rtc.ALARM0, 300000) # 5 * 60 * 1000ms
                machine.deepsleep()
    except OSError as e:
        print("check_msg:", e)
        print("try to reconnect")
        print("deepsleep for 60 sec and retry")
        rtc.alarm(rtc.ALARM0, 60000) # 60 * 1000ms
        machine.deepsleep()

print("ttime = ", ttime)

# first, check if a cached 'cfeed' exists:
try:
    print("Looking for cached 'cfeed'")
    f = open('holdover.txt', 'r') # f.write( strtosend)
    print("Found cached 'cfeed', reading and trying to send")
    # split between topic and message
    strtosend = f.read()
    try:
        c.publish(b"cfeed", strtosend)
    except OSError as e:
        print("can't publish, deep sleeping for 1 min and will try again")
        rtc.alarm(rtc.ALARM0, 60000) # 60 * 1000ms
        machine.deepsleep()
    f.close()
    os.remove('holdover.txt')
except OSError as e:
    print("no cached 'cfeed', continuing")

# Get orders
orders = str()
while orders == '':
    print("New orders being retrieved")
    c.connect(clean_session=True)

    print('Retrieving from master: ' + master_topic)
    c.subscribe(master_topic)
    print('Retrieving from specific: ' + sub_topic)
    c.subscribe(sub_topic)
    while orders == '':
        c.check_msg()
        time.sleep(1)

global nextfeed
## Now that we have our ttime, marching orders, and calibration, we can check these and feed:
need2feed = False
while ttime != (DEFAULT_TIME):
    try:
        f = open("need2feed", "r")
        need2feed = True
        f.close()
    except:
        need2feed = False

    if need2feed:
        os.remove("need2feed")
        # Feed cycle
        ## Conditioning
        lcond.on()
        time.sleep(conditioning_time)
        ## Feeding
        servo_pos(servo, 0 + food_hole_offset)
        time.sleep(1) # let it feed
        servo_pos(servo, 0)
        lcond.off()
        ## Publish response
        response = CLIENT_ID + ',' + str(ttime)
        try:
            c.publish(b"cfeed", response)
        except OSError as e:
            print("Not able to pub, saving file then deepsleeping for 30 seconds")
            # need to write and save for retrying later
            f = open("holdover.txt", 'w')
            f.write(strtosend)
            f.close()
            time.sleep(3) # Wait for fs to complete
            rtc.alarm(rtc.ALARM0, 30000)
            machine.deepsleep() # deepsleep for 30s
        print('waiting')
        time.sleep(3)
        rtc.alarm(rtc.ALARM0, postfeedsleepms)
        print('Deep sleeping for postfeed: ', postfeedsleepms)
        machine.deepsleep()

    if orders > ttime:
        nextfeed = orders - ttime
    else:
        nextfeed = default_waiting_time

    # put into deep sleep for 80% of the time to the next feed
    print('Next feeding is in ' + str(nextfeed))
    if nextfeed > 120:
        sleepfor = int(nextfeed * 0.8 * 1000) # milliseconds
        print('Deep sleeping for (s): ' + str(sleepfor / 1000))
        rtc.alarm(rtc.ALARM0, sleepfor)
        machine.deepsleep()

    if (120) >= (nextfeed) > 0:
        sleepfor = int(nextfeed * 1000)
        print('Deep sleeping briefly for (s): ' + str(nextfeed / 1000))
        f = open("need2feed", "w")
        f.close()
        rtc.alarm(rtc.ALARM0, sleepfor)
        machine.deepsleep()
