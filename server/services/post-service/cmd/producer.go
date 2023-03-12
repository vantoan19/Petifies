package cmd

import (
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var PostProducer producer.KafkaProducer

func initPostProducer() error {
	postProducer, err := producer.NewKafkaProducer(config.NewKafkaProducerConfig(Conf.Brokers), logging.New("PostService.Producer"))
	if err != nil {
		return err
	}
	PostProducer = *postProducer
	return nil
}
