
#include <MQTT.h>
#include <ESP8266mDNS.h>
#include <ESP8266WiFi.h>
#include <Servo.h>

extern "C" {
#include "user_interface.h" // this is for the RTC memory read/write functions
}

#define WIFI_NAME "Samarinda"
#define WIFI_PASSWORD "whatisyourname"

#define MASTER_TOPIC "orders/all"
#define SPECIFIC_TOPIC "orders/feeder"

#define DEFAULT_WAITING_TIME 3600 // 1 hour in seconds
#define RTCMEMORYSTART 65

typedef struct {
  int nextfeed;
  bool use_master;
  bool execute_feed;
} rtcStore;

rtcStore rtcMem;

int default_pos = 90;
int food_hole_offset = 30; // Angles

int conditioningLED = 14;
int servoPin = 2;

bool confirm_time;
bool confirm_feed;
bool check_feed;

Servo servo;
WiFiClient net;
MQTTClient mqtt;
char id[34]; // F + uint32 (hex) + null byte 

int ttime;
int nextfeed;

void setup() {
  ttime = -1;
  nextfeed = 0;
  
  confirm_time = false;
  confirm_feed = false;
  check_feed = false;
  
  Serial.begin(115200);
  while(!Serial) { }

  WiFi.begin(WIFI_NAME, WIFI_PASSWORD);
  Serial.print("Connecting to WiFi");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
    Serial.flush();
  }
  Serial.println();

  sprintf(id, "F%08x", ESP.getChipId());

  MDNS.begin(id);
  int n = MDNS.queryService("mqtt", "tcp");
  Serial.println(n);
  if (n == 0) {
    Serial.println("cannot find MQTT over mDNS, sleeping for 5 minutes");
    ESP.deepSleep(5 * 60 * 1000 * 1000); // microseconds
  } else {
    mqtt.begin(MDNS.IP(0), net);
    mqtt.onMessage(messageReceived);

    Serial.print("Connecting to MQTT at ");
    Serial.print(MDNS.IP(0));
    while(!mqtt.connect(id, "user", "password")) {
      Serial.print(".");
      Serial.flush();
      delay(500);
    }
  
    Serial.println("\nConnected to MQTT!");
    mqtt.subscribe("time");
    mqtt.subscribe(MASTER_TOPIC);
    mqtt.subscribe(SPECIFIC_TOPIC);
  }

  pinMode(conditioningLED, OUTPUT);
  digitalWrite(conditioningLED, LOW);

  servo.attach(servoPin);
  servo.write(default_pos);

  rtcMem.use_master = true;
  system_rtc_mem_read(RTCMEMORYSTART, &rtcMem, sizeof(rtcMem));
  Serial.println(rtcMem.nextfeed);

  if (rtcMem.execute_feed) {
    feed();
    
    confirm_feed = true;

    rtcMem.execute_feed = false;
    system_rtc_mem_write(RTCMEMORYSTART, &rtcMem, sizeof(rtcMem));
  }
}

void loop() {
  mqtt.loop();
  if (confirm_time) {
    confirm_time = false;
    
    char res[128];
    sprintf(res, "%s,%d", id, ttime);
    mqtt.publish("ctime", res);
  }

  if (confirm_feed) {
    if (ttime != -1) {
      confirm_feed = false;
      
      char res[128];
      sprintf(res, "%s,%d", id, ttime);
      mqtt.publish("cfeed", res);
    }
  }

  if (check_feed) {
    if (ttime != -1) {
      int sleep_for = rtcMem.nextfeed > ttime ? rtcMem.nextfeed - ttime : DEFAULT_WAITING_TIME;
      if (sleep_for > 120) {
        sleep_for = sleep_for * 0.8 * 1000 * 1000; // 80%, then convert to microseconds
        Serial.print("Sleeping for (s): ");
        Serial.println(sleep_for / 1000000);
      } else {
        sleep_for = sleep_for * 1000 * 1000; // Converting to microseconds
        Serial.println("Sleeping until next feed");

        rtcMem.execute_feed = true;
        system_rtc_mem_write(RTCMEMORYSTART, &rtcMem, sizeof(rtcMem));
      }

      check_feed = false;
      ESP.deepSleep(sleep_for); // microseconds
    }
  }
  
  if (Serial.available() > 0) {
    char ch = Serial.read();
    if (ch == 't') {
      Serial.println(ttime);
    } else if (ch == 's') {
      Serial.println("sleeping");
      Serial.println(ttime);
      Serial.println("---");
  
      rtcMem.nextfeed = ttime;
      system_rtc_mem_write(RTCMEMORYSTART, &rtcMem, sizeof(rtcMem));
      ESP.deepSleep(3 * 1000 * 1000); // microseconds
    }
  }
}

void messageReceived(String &topic, String &payload) {
  if (topic == "time") {
    ttime = payload.toInt();
    confirm_time = true;
  } else if (topic == MASTER_TOPIC) {
    if (rtcMem.use_master) {
      process_feed_order(payload); 
    }
  } else if (topic == SPECIFIC_TOPIC) {
    Serial.println("received");
    process_feed_order(payload);
    rtcMem.use_master = false;
    system_rtc_mem_write(RTCMEMORYSTART, &rtcMem, sizeof(rtcMem));
  }
}

void process_feed_order(String &raw_feed) {
  rtcMem.nextfeed = raw_feed.toInt();
  system_rtc_mem_write(RTCMEMORYSTART, &rtcMem, sizeof(rtcMem));
  check_feed = true;
}

void feed() {
  digitalWrite(conditioningLED, HIGH);
  delay(3000);
  servo.write(default_pos + food_hole_offset * 2);
  delay(1000);
  servo.write(default_pos); 
  digitalWrite(conditioningLED, LOW);
}
