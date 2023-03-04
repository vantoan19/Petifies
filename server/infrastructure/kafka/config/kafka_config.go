package config

import "github.com/Shopify/sarama"

type KafkaConfig struct {
	ProducerConfig *sarama.Config
	ConsumerConfig *sarama.Config
	Brokers        []string
	Topic          string
	ConsumerGroup  string
}

func NewKafkaConfig(brokers []string, topic string, consumerGroup string) *KafkaConfig {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true

	return &KafkaConfig{
		ProducerConfig: producerConfig,
		ConsumerConfig: consumerConfig,
		Brokers:        brokers,
		Topic:          topic,
		ConsumerGroup:  consumerGroup,
	}
}
