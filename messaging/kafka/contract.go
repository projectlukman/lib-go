// Package kafka
package kafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
)

// Producer represents kafka publisher message topic
type Producer interface {
	Publish(ctx context.Context, msg *MessageContext) error
}

// Consumer represents a Sarama consumer consumer interface
type Consumer interface {
	Subscribe(*ConsumerContext)
}

type MessageContext struct {
	Value     string
	Key       []byte
	LogId     interface{}
	Topic     string
	Partition int32
	Offset    int64
	TimeStamp time.Time
	Verbose   bool
}

type ConsumerContext struct {
	Handler MessageProcessorFunc
	Topics  []string
	GroupID string
	Context context.Context
}

var balanceStrategies = map[string]sarama.BalanceStrategy{
	sarama.RoundRobinBalanceStrategyName: sarama.BalanceStrategyRoundRobin,
	sarama.RangeBalanceStrategyName:      sarama.BalanceStrategyRange,
	sarama.StickyBalanceStrategyName:     sarama.BalanceStrategySticky,
}
