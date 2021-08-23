package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Options struct {
	StoreDir					string
	Offset						time.Duration
	MQTT							MQTTClientOptions
	DefaultTag				string
	DefaultSchedule		string
}

type Handler struct {
	db					*DBClient
	scheduler		*Scheduler
	mqtt				*MQTTClient
	template		*template.Template
}

func NewHandler(options Options) (*Handler, error) {
	db, err := NewDBClient(options.StoreDir)
	if err != nil {
		return nil, err
	}

	mqtt, err := NewMQTT(options.MQTT)
	if err != nil {
		return nil, err
	}

	template, err := template.ParseFiles("./pages/index.html")
	if err != nil {
		return nil, err
	}

	scheduler := NewScheduler(options.Offset)

	h := &Handler {
		db,
		scheduler,
		mqtt,
		template,
	}

	jobs, err := h.db.ListKeys()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		err := h.scheduler.Add(job.Tag, job.Cron, h.RunJob)
		if err != nil {
			return nil, err
		}
	}

	if len(jobs) == 0 {
		log.Println("no jobs cached, adding default job")
		err := h.scheduler.Add(options.DefaultTag, options.DefaultSchedule, h.RunJob)
		if err != nil {
			return nil, err
		}

		err = h.db.InsertEntry(options.DefaultTag, options.DefaultSchedule)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

func (h *Handler) Close() error {
	h.scheduler.Stop()
	return h.db.Close()
}

func (h *Handler) Error(w http.ResponseWriter, r *http.Request, code int, msg string) {
	accepts := NewTypes(r, "Accept")
	if accepts.Has("text/html") {
		ctx := context.WithValue(r.Context(), "error-code", code)
		ctx = context.WithValue(ctx, "error-msg", msg)
		http.Redirect(w, r.WithContext(ctx), "/", http.StatusSeeOther)
	} else {
		http.Error(w, msg, code)
	}
}

type IndexData struct {
	Error string
	Orders []RawScheduleRequest
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accepts := NewTypes(r, "Accept")

	orders, err := h.db.ListKeys()
	if err != nil && !accepts.Has("text/html") {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if accepts.Has("text/html") {
		err2 := ""
		raw_error_msg := r.Context().Value("error-msg")
		if raw_error_msg != nil {
			err2 = raw_error_msg.(string)
		}

		if err != nil {
			err2 = err2 + "; " + err.Error()
		}

		data := IndexData {
			err2,
			orders,
		}

		h.template.Execute(w, &data)
	} else if accepts.Has("application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&orders)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		for _, order := range orders {
			w.Write([]byte(order.Tag + " " + order.Cron))
		}
	}
}

type RawScheduleRequest struct {
	Tag		string	`json:"tag"`
	Cron	string	`json:"schedule"`
}

func (h *Handler) PostScheduleRaw(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	contentTypes := NewTypes(r, "Content-Type")

	var data RawScheduleRequest
	if contentTypes.Has("application/json") {
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			h.Error(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}
	} else if contentTypes.Has("application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			h.Error(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		tag := r.Form.Get("tag")
		if tag == "" {
			h.Error(w, r, http.StatusUnprocessableEntity, "field 'tag' is missing")
			return
		}

		schedule := r.Form.Get("schedule")
		if schedule == "" {
			h.Error(w, r, http.StatusUnprocessableEntity, "field 'schedule' is missing")
			return
		}

		data = RawScheduleRequest {
			tag,
			schedule,
		}
	} else {
		h.Error(w, r, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
		return
	}

	err := h.scheduler.Add(data.Tag, data.Cron, h.RunJob)
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.db.InsertEntry(data.Tag, data.Cron)
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	accepts := NewTypes(r, "Accept")
	if accepts.Has("text/html") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tag := strings.TrimPrefix(p.ByName("tag"), "/")
	h.scheduler.Remove(tag)

	err := h.db.RemoveEntry(tag)
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	accepts := NewTypes(r, "Accept")
	if accepts.Has("text/html") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) RunJob(j *Job) error {
	err := h.mqtt.Publish(j.Tag, j.NextFeed)
	if err != nil {
		return err
	}

	return nil
}
