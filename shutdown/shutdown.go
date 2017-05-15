package shutdown

import (
	"github.com/pufferpanel/pufferd/programs"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/braintree/manners"
	"os/signal"
	"os"
	"runtime/debug"
	"sync"
)

func CompleteShutdown() {
	logging.Warn("Interrupt received, stopping servers")
	wg := Shutdown()
	wg.Wait()
	logging.Info("All servers stopped")
	os.Exit(0)
}

func Shutdown() *sync.WaitGroup {
	defer func() {
		if err := recover(); err != nil {
			logging.Errorf("Error: %+v\n%s", err, debug.Stack())
		}
	}()
	wg := sync.WaitGroup{}
	manners.Close()
	prgs := programs.GetAll()
	wg.Add(len(prgs))
	for _, element := range prgs {
		go func(e programs.Program) {
			defer func() {
				if err := recover(); err != nil {
					logging.Errorf("Error: %+v\n%s", err, debug.Stack())
				}
				wg.Done()
			}()
			logging.Warn("Stopping program " + e.Id())
			err := e.Stop()
			if err != nil {
				logging.Error("Error stopping server " + e.Id(), err)
			}
			err = e.GetEnvironment().WaitForMainProcess()
			if err != nil {
				logging.Error("Error stopping server " + e.Id(), err)
			}
			logging.Warn("Stopped program " + e.Id())
		}(element)
	}
	return &wg;
}

func CreateHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logging.Errorf("Error: %+v\n%s", err, debug.Stack())
			}
		}()
		<- c
		CompleteShutdown()
	}()
}