package consumer

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       sarama.ConsumerGroupHandler
	topic         string
}

type KafkaMessageHandler func(ctx context.Context, message *models.KafkaMessage) error

type kafkaConsumerHandler struct {
	handler KafkaMessageHandler
}

func NewKafkaConsumer(kafkaConfig *config.KafkaConsumerConfig, handler KafkaMessageHandler) (*KafkaConsumer, error) {
	consumerGroup, err := sarama.NewConsumerGroup(kafkaConfig.Brokers, kafkaConfig.ConsumerGroup, kafkaConfig.ConsumerConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		handler: &kafkaConsumerHandler{
			handler: handler,
		},
		topic: kafkaConfig.Topic,
	}, nil
}

func (c *KafkaConsumer) Consume() {
	ctx := context.Background()
	for {
		err := c.consumerGroup.Consume(ctx, []string{c.topic}, c.handler)
		if err != nil {
			log.Printf("Error while consuming: %v", err)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.consumerGroup.Close()
}

func (h *kafkaConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *kafkaConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *kafkaConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.Background()
	for {
		select {
		case msg := <-claim.Messages():
			err := h.handler(ctx, &models.KafkaMessage{
				Topic:     msg.Topic,
				Partition: msg.Partition,
				Offset:    msg.Offset,
				Key:       msg.Key,
				Value:     msg.Value,
			})
			if err != nil {
				return err
			}
			session.MarkMessage(msg, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
