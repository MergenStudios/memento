package main

import (
	"fmt"
	"github.com/kardianos/service"
	"memento/scripts"
	"sync"
)

var (
	serviceIsRunning bool
	programIsRunning bool
	writingSync      sync.Mutex
)

const serviceName = "Memento sources watcher"
const serviceDescription = "Background service for memento watching and updating permanent data sources (Powermachine#1688 message me on discord if you found this)"

type program struct{}

func (p program) Start(s service.Service) error {
	writingSync.Lock()
	serviceIsRunning = true
	writingSync.Unlock()

	// update the permanent sources here
	scripts.Update("all-background")
	return nil
}

func (p program) Stop(s service.Service) error {
	return nil
}

func (p program) run() {
	select {}
}

func main() {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println("Cannot create the service: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		fmt.Println("Cannot start the service: " + err.Error())
	}
}
