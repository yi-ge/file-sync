package main

import (
	"time"

	"github.com/kardianos/service"
)

// Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		// logger.Info("Running in terminal.")

		command(s)
	} else {
		logger.Info("Running under service manager.")

		// Start should not block. Do the actual work async.
		go p.run()
	}
	p.exit = make(chan struct{})

	return nil
}

func (p *program) run() error {
	data, err := getData()

	if err == nil {
		go StartSSEClient(data)
	} else {
		logger.Error(err)
	}

	logger.Infof("I'm running %v.", service.Platform())
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case tm := <-ticker.C:
			logger.Infof("Still running at %v...", tm)
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}
