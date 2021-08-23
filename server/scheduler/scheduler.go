package main

import (
	"log"
	"sort"
	"time"

	"github.com/robfig/cron/v3"
)

type Job struct {
	Tag				string
	Offset		time.Duration
	Schedule	cron.Schedule
	Fn				func(*Job) error

	NextFeed	time.Time
	NextRun		time.Time
}

func (j *Job) SetNext(t time.Time) {
	j.NextFeed = j.Schedule.Next(t.Add(j.Offset))
	j.NextRun = j.NextFeed.Add(-1 * j.Offset)
}

func (j *Job) Run() {
	err := j.Fn(j)
	if err != nil {
		log.Println(err)
	}
}

type Jobs []*Job
func (j Jobs) Len() int {
	return len(j)
}

func (x Jobs) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x Jobs) Less(i, j int) bool {
	if x[i].NextRun.IsZero() {
		return false
	}

	if x[j].NextRun.IsZero() {
		return true
	}

	return x[i].NextRun.Before(x[j].NextRun)
}

func (j Jobs) Has(tag string) bool {
	for _, job := range j {
		if job.Tag == tag {
			return true
		}
	}

	return false
}

type Scheduler struct {
	offset		time.Duration
	parser		cron.Parser
	jobs			Jobs
	running		bool

	add				chan *Job
	remove		chan string
	stop			chan struct{}
}

func NewScheduler(offset time.Duration) *Scheduler {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	jobs := make(Jobs, 0)

	add := make(chan *Job)
	remove := make(chan string)
	stop := make(chan struct{})

	s := &Scheduler {
		offset,
		parser,
		jobs,
		false,
		add,
		remove,
		stop,
	}

	go s.run()

	return s
}

func (s *Scheduler) Stop() {
	if !s.running {
		return
	}

	s.running = false
	s.stop <- struct{}{}
}

func (s *Scheduler) run() {
	if s.running {
		return
	}

	s.running = true
	now := time.Now()
	for _, job := range s.jobs {
		job.SetNext(now)
	}

	for {
		sort.Sort(s.jobs)
		now = time.Now()

		var timer *time.Timer
		if len(s.jobs) == 0 || s.jobs[0].NextRun.IsZero() {
			timer = time.NewTimer(100000 * time.Hour) // Some absurd duration
		} else {
			timer = time.NewTimer(s.jobs[0].NextRun.Sub(now))
		}

		for {
			select {
				case n := <-timer.C:
					for _, j := range s.jobs {
						if j.NextRun.After(n) || j.NextRun.IsZero() {
							break
						}

						go j.Run()
						j.SetNext(n)
					}

				case e := <-s.add:
					timer.Stop()
					now = time.Now()
					e.SetNext(now)
					if !s.jobs.Has(e.Tag) {
						s.jobs = append(s.jobs, e)
					} else {
						for _, j := range s.jobs {
							if j.Tag == e.Tag {
								j.Schedule = e.Schedule
								j.Offset = e.Offset
								j.Fn = e.Fn
								j.NextRun = e.NextRun
								break
							}
						}
					}

				case tag := <-s.remove:
					timer.Stop()
					for i, job := range s.jobs {
						if job.Tag == tag {
							s.jobs = append(s.jobs[:i], s.jobs[i+1:]...)
							break
						}
					}

				case <-s.stop:
					timer.Stop()
					return
			}

			break
		}
	}
}


func (s *Scheduler) Add(tag string, cronstring string, fn func(j *Job) error) error {
	job, err := s.NewJob(tag, cronstring, fn)
	if err != nil {
		return err
	}

	if !s.running {
		if !s.jobs.Has(job.Tag) {
			s.jobs = append(s.jobs, job)
		} else {
			for _, j := range s.jobs {
				if j.Tag == job.Tag {
					j.Schedule = job.Schedule
					j.Offset = job.Offset
					j.Fn = job.Fn
					j.NextRun = job.NextRun
					break
				}
			}
		}
	} else {
		s.add <- job
	}

	return nil
}

func (s *Scheduler) Remove(tag string) {
	if !s.running {
		for i, job := range s.jobs {
			if job.Tag == tag {
				s.jobs = append(s.jobs[:i], s.jobs[i+1:]...)
				break
			}
		}
	} else {
		s.remove <- tag
	}
}

func (s *Scheduler) NewJob(tag string, cronstring string, fn func(j *Job) error) (*Job, error) {
	schedule, err := s.parser.Parse(cronstring)
	if err != nil {
		return nil, err
	}

	j := &Job {
		Tag: tag,
		Offset: s.offset,
		Schedule: schedule,
		Fn: fn,
	}

	return j, nil
}
