package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gookit/ini"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var mqttConf string
	flag.StringVar(&mqttConf, "mqtt", "../server.conf", "path to server.conf")
	var schedulerConf string
	flag.StringVar(&schedulerConf, "scheduler", "./scheduler.conf", "path to scheduler.conf")
	flag.Parse()

	serverConf, err := ini.LoadExists(mqttConf, schedulerConf)
	if err != nil {
		log.Fatalln(err)
		return
	}

	mqttHost, ok := serverConf.String("MQTT_HOST")
	if !ok {
		mqttHost = "127.0.0.1"
	}

	mqttPort, ok := serverConf.String("MQTT_PORT")
	if !ok {
		mqttPort = "1883"
	}

	mqttUsername, ok := serverConf.String("MQTT_USER")
	if !ok {
		mqttUsername = "user"
	}

	mqttPassword, ok := serverConf.String("MQTT_PASS")
	if !ok {
		mqttPassword = "password"
	}

	clientID, ok := serverConf.String("CLIENT_ID")
	if !ok {
		clientID = "scheduler"
	}

	mqttAddress := mqttHost + ":" + mqttPort

	mqttClientOptions := MQTTClientOptions {
		Address: mqttAddress,
		ClientID: clientID,
		Username: mqttUsername,
		Password: mqttPassword,
	}

	storeDir, ok := serverConf.String("STORE_DIR")
	if !ok {
		storeDir = "./tmp"
	}

	defaultTag, ok := serverConf.String("DEFAULT_TAG")
	if !ok {
		defaultTag = "orders/all"
	}

	defaultSchedule, ok := serverConf.String("DEFAULT_SCHEDULE")
	if !ok {
		defaultSchedule = "0 16 * * *"
	}

	options := Options {
		StoreDir: storeDir,
		Offset: time.Second,
		MQTT: mqttClientOptions,
		DefaultTag: defaultTag,
		DefaultSchedule: defaultSchedule,
	}

	handler, err := NewHandler(options)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer handler.Close()

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("./static"))
	router.GET("/", handler.Index)
	router.POST("/schedule", handler.PostScheduleRaw)
	router.GET("/schedule/remove/*tag", handler.Remove)

	listen, _ := serverConf.String("LISTEN")
	log.Printf("listening on %s\n", listen)
	log.Fatalln(http.ListenAndServe(listen, router))
}
