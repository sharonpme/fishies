package main

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	mqttClientOptions := MQTTClientOptions {
		Address: "172.16.12.101:1883",
		ClientID: "scheduler",
		Username: "user",
		Password: "password",
	}

	options := Options {
		StoreDir: "./tmp",
		Offset: time.Second,
		MQTT: mqttClientOptions,
		DefaultTag: "orders/all",
		DefaultSchedule: "1 * * * *",
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

	listen := ":1837"
	log.Fatalln(http.ListenAndServe(listen, router))
}
