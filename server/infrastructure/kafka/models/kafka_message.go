package models

type KafkaModel interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

type KafkaMessage struct {
	Topic     string     `json:"topic"`
	Partition int32      `json:"partition"`
	Offset    int64      `json:"offset"`
	Key       []byte     `json:"key"`
	Value     KafkaModel `json:"value"`
}
