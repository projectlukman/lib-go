// Package kafka messaging
// @author Daud Valentino
package kafka

import (
	"github.com/Shopify/sarama"
)

// Consumer represents a Sarama consumer group consumer
type consumerHandler struct {
	msgProcessor MessageProcessorFunc
	autoCommit   bool
}

// NewConsumerHandler return consumer handler
func NewConsumerHandler(msgProcessor MessageProcessorFunc, autoCommit bool) sarama.ConsumerGroupHandler {
	return &consumerHandler{msgProcessor: msgProcessor, autoCommit: autoCommit}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for msg := range claim.Messages() {
		if c.autoCommit {
			session.MarkMessage(msg, "")
		}
		c.msgProcessor(&MessageDecoder{
			Body:      msg.Value,
			Key:       msg.Key,
			Partition: msg.Partition,
			TimeStamp: msg.Timestamp,
			Offset:    msg.Offset,
			Topic:     msg.Topic,
			Commit: func(m *MessageDecoder) {
				session.MarkOffset(m.Topic, m.Partition, m.Offset+1, "")
			},
		})
	}

	return nil
}
