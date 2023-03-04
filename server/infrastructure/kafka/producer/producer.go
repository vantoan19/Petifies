package producer

import (
	"github.com/Shopify/sarama"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(kafkaConfig *config.KafkaProducerConfig) (*KafkaProducer, error) {
	producer, err := sarama.NewSyncProducer(kafkaConfig.Brokers, kafkaConfig.ProducerConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
	}, nil
}

func (p *KafkaProducer) SendMessage(msg *models.KafkaMessage) (*models.KafkaMessage, error) {
	bytes, err := msg.Value.Serialize()
	if err != nil {
		return nil, err
	}

	pmsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(bytes),
	}
	partition, offset, err := p.producer.SendMessage(pmsg)
	if err != nil {
		return nil, err
	}
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
