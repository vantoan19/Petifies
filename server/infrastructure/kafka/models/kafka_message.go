package models

type KafkaModel interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

type KafkaMessage struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     KafkaModel
}
