package outbox_dispatcher

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

type DispatcherSettings struct {
	PublishInterval time.Duration
	UnlockInterval  time.Duration
	CleanInterval   time.Duration
	PublishSettings
	CleanSettings
	UnlockSettings
}

type Dispatcher interface {
	Run(errsChan chan<- error, endSignal <-chan bool)
}

type dispatcher struct {
	settings  DispatcherSettings
	unlocker  unlocker
	cleaner   cleaner
	publisher publisher
	logger    logging.Logger
}

func NewDispatcher(repo outbox_repo.EventRepository, producer producer.KafkaProducer, settings DispatcherSettings, logger logging.Logger) Dispatcher {
	return &dispatcher{
		settings:  settings,
		unlocker:  NewUnlocker(repo, settings.UnlockSettings),
		cleaner:   NewCleaner(repo, settings.CleanSettings),
		publisher: NewPublisher(producer, repo, settings.PublishSettings),
		logger:    logger,
	}
}

func (d *dispatcher) Run(errsChan chan<- error, endSignal <-chan bool) {
	d.logger.Info("Start EventDispatcher")

	endUnlocker := make(chan bool, 1)
	endCleaner := make(chan bool, 1)
	endPublisher := make(chan bool, 1)

	go d.runUnlocker(errsChan, endUnlocker)
	go d.runCleaner(errsChan, endCleaner)
	go d.runPublisher(errsChan, endPublisher)

	go func() {
		// block until end signal arrives
		<-endSignal
		d.logger.Info("Received end signal, stopping the event dispatcher")
		endUnlocker <- true
		endCleaner <- true
		endPublisher <- true
		d.logger.Info("Finish EventDispatcher")
	}()
}

func (d *dispatcher) runUnlocker(errsChan chan<- error, endSignal <-chan bool) {
	ticker := time.NewTicker(d.settings.UnlockInterval)
	for {
		d.logger.Info("Unlocking long-time locked events")
		err := d.unlocker.UnlockLongtimeLockedEvents()
		if err != nil {
			d.logger.ErrorData("Error at event unlocker", logging.Data{"error": err.Error()})
			errsChan <- err
		}

		select {
		case <-ticker.C:
			continue
		case <-endSignal:
			d.logger.Info("Received end signal, stopping unlocker")
			ticker.Stop()
			return
		}
	}
}

func (d *dispatcher) runCleaner(errsChan chan<- error, endSignal <-chan bool) {
	ticker := time.NewTicker(d.settings.CleanInterval)
	for {
		d.logger.Info("Cleaning expired events")
		err := d.cleaner.DeleteExpiredEvents()
		if err != nil {
			d.logger.ErrorData("Error at event cleaner", logging.Data{"error": err.Error()})
			errsChan <- err
		}

		select {
		case <-ticker.C:
			continue
		case <-endSignal:
			d.logger.Info("Received end signal, stopping cleaner")
			ticker.Stop()
			return
		}
	}
}

func (d *dispatcher) runPublisher(errsChan chan<- error, endSignal <-chan bool) {
	ticker := time.NewTicker(d.settings.PublishInterval)
	for {
		d.logger.Info("Publishing events")
		err := d.publisher.PublishEvents()
		if err != nil {
			d.logger.ErrorData("Error at event publisher", logging.Data{"error": err.Error()})
			errsChan <- err
		}

		select {
		case <-ticker.C:
			continue
		case <-endSignal:
			d.logger.Info("Received end signal, stopping publisher")
			ticker.Stop()
			return
		}
	}
}
