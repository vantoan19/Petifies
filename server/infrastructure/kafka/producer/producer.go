package producer

import (
	"github.com/Shopify/sarama"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(kafkaConfig *config.KafkaConfig) (*KafkaProducer, error) {
	producer, err := sarama.NewSyncProducer(kafkaConfig.Brokers, kafkaConfig.ProducerConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    kafkaConfig.Topic,
	}, nil
}

func (p *KafkaProducer) SendMessage(key []byte, value models.KafkaModel) (*models.KafkaMessage, error) {
	bytes, err := value.Serialize()
	if err != nil {
		return nil, err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(bytes),
	}
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return nil, err
	}
	return &models.KafkaMessage{
		Topic:     p.topic,
		Partition: partition,
		Offset:    offset,
		Key:       key,
		Value:     value,
	}, nil
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
