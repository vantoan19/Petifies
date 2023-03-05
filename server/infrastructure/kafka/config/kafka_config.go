package config

import "github.com/Shopify/sarama"

type KafkaProducerConfig struct {
	ProducerConfig *sarama.Config
	Brokers        []string
}

func NewKafkaProducerConfig(brokers []string) *KafkaProducerConfig {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true

	return &KafkaProducerConfig{
		ProducerConfig: producerConfig,
		Brokers:        brokers,
	}
}

type KafkaConsumerConfig struct {
	ConsumerConfig *sarama.Config
	Brokers        []string
	Topic          string
	ConsumerGroup  string
}

func NewKafkaConsumerConfig(brokers []string, topic string, consumerGroup string) *KafkaConsumerConfig {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true

	return &KafkaConsumerConfig{
		ConsumerConfig: consumerConfig,
		Brokers:        brokers,
		Topic:          topic,
		ConsumerGroup:  consumerGroup,
	}
}
