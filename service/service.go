package service

import (
	"log"
	"math"
	"time"

	"github.com/boltdb/bolt"
	"github.com/e-dard/lumen/integration"
)

type Service struct {
	db      *bolt.DB
	quit    chan struct{}
	ticker  *time.Ticker
	verbose bool
	tasks   map[string]task
}

type task struct {
	Name   []byte
	Bounds [2]int64
	Task   integration.Runnable
}

type Option func(*Service)

func WithInterval(d time.Duration) func(*Service) {
	return func(s *Service) {
		s.ticker = time.NewTicker(d)
	}
}

func Verbose() func(*Service) {
	return func(s *Service) {
		s.verbose = true
	}
}

func NewService(dbpth string, options ...Option) *Service {
	db, err := bolt.Open(dbpth, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	s := &Service{
		db:     db,
		quit:   make(chan struct{}),
		ticker: time.NewTicker(10 * time.Second),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *Service) Close() {
	s.quit <- struct{}{}
	s.ticker.Stop()
	log.Println("Service closed OK.")
}

func (s *Service) Start() error {
	if err := s.loadTasks(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.update()
			case <-s.quit:
				return
			}
		}
	}()
	return nil
}

func (s *Service) loadTasks() error {
	s.tasks = map[string]task{
		"st-dark": task{
			Name:   []byte("st-dark"),
			Bounds: [2]int64{0, 32700000},
			Task: integration.NewSublimeTask(
				"/Users/edd/Dropbox/app_options/ST3/Packages/User/Preferences.sublime-settings",
				[]byte(`{"color_scheme": "Packages/Color Scheme - Default/Solarized (Dark).tmTheme"}`),
			),
		},
		"st-light": task{
			Name:   []byte("st-light"),
			Bounds: [2]int64{32700000, math.MaxInt64},
			Task: integration.NewSublimeTask(
				"/Users/edd/Dropbox/app_options/ST3/Packages/User/Preferences.sublime-settings",
				[]byte(`{"color_scheme": "Packages/Color Scheme - Default/Solarized (Light).tmTheme"}`),
			),
		},
	}
	return nil
}

func (s *Service) debugf(format string, v ...interface{}) {
	if s.verbose {
		log.Printf(format, v...)
	}
}

func (s *Service) debug(v ...interface{}) {
	if s.verbose {
		log.Println(v...)
	}
}

func ReadSensors() (Reading, error) {
	reading, err := readSensors()
	if err == ErrNoLightSensors {
		log.Fatal(err)
	}
	return reading, err
}

func (s *Service) update() {
	// read sensors
	reading, err := ReadSensors()
	if err != nil {
		log.Println(err)
	}
	s.debug(reading)

	// Run tasks?
	for _, t := range s.tasks {
		if reading.Mean >= t.Bounds[0] && reading.Mean < t.Bounds[1] {
			s.debug("running", string(t.Name))
			if err := t.Task.Run(); err != nil {
				log.Println(err)
			}
		}
	}
}
