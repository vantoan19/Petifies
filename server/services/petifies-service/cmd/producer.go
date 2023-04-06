package cmd

import (
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var PetifiesProducer producer.KafkaProducer

func initPostProducer() error {
	petifiesProducer, err := producer.NewKafkaProducer(config.NewKafkaProducerConfig(Conf.Brokers), logging.New("PetifiesService.Producer"))
	if err != nil {
		return err
	}
	PetifiesProducer = *petifiesProducer
	return nil
}
