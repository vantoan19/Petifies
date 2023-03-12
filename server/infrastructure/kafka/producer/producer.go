package producer

import (
	"github.com/Shopify/sarama"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	logger   *logging.Logger
}

func NewKafkaProducer(kafkaConfig *config.KafkaProducerConfig, logger *logging.Logger) (*KafkaProducer, error) {
	producer, err := sarama.NewSyncProducer(kafkaConfig.Brokers, kafkaConfig.ProducerConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		logger:   logger,
	}, nil
}

func (p *KafkaProducer) SendMessage(msg *models.KafkaMessage) (*models.KafkaMessage, error) {
	p.logger.InfoData("Publishing Event", logging.Data{"topic": msg.Topic, "key": msg.Key})

	pmsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	}
	partition, offset, err := p.producer.SendMessage(pmsg)
	if err != nil {
		p.logger.ErrorData("Error at publishing event", logging.Data{"topic": msg.Topic, "key": msg.Key, "error": err.Error()})
		return nil, err
	}

	p.logger.InfoData("Published Event", logging.Data{"topic": msg.Topic, "key": msg.Key})
	return &models.KafkaMessage{
		Topic:     msg.Topic,
		Partition: partition,
		Offset:    offset,
		Key:       msg.Key,
		Value:     msg.Value,
	}, nil
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
