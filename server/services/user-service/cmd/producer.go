package cmd

import (
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var UserProducer producer.KafkaProducer

func initUserProducer() error {
	userProducer, err := producer.NewKafkaProducer(config.NewKafkaProducerConfig(Conf.Brokers), logging.New("UserService.Producer"))
	if err != nil {
		return err
	}
	UserProducer = *userProducer
	return nil
}
