// Package kafka messaging broker
// @author Daud Valentino
package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/projectlukman/lib-go/log"
)

const (
	defaultTimeout = 3 // in second
)

var (
	partitiions = map[string]sarama.PartitionerConstructor{
		"hash":       sarama.NewHashPartitioner,
		"roundrobin": sarama.NewRoundRobinPartitioner,
		"reference":  sarama.NewReferenceHashPartitioner,
		"random":     sarama.NewRandomPartitioner,
		"manual":     sarama.NewManualPartitioner,
	}
)

type producer struct {
	config   *sarama.Config
	brokers  []string
	producer sarama.SyncProducer
}

// SyncPublisher publish message  synchronously
func (k *producer) Publish(_ context.Context, msg *MessageContext) error {

	key, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("fail to create producer key")
	}
	param := &sarama.ProducerMessage{
		Topic:     msg.Topic,
		Value:     sarama.StringEncoder(msg.Value),
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Timestamp: msg.TimeStamp,
		Key:       sarama.StringEncoder(key.String()),
	}

	if msg.Key != nil && len(msg.Key) > 0 {
		param.Key = sarama.ByteEncoder(msg.Key)
	}

	partition, offset, err := k.producer.SendMessage(param)

	if err != nil {
		return fmt.Errorf("[kafka-publisher] topic: %s, partition %d, offset %d, id %v, got: %w", msg.Topic, partition, offset, msg.LogId, err)
	}

	if msg.Verbose {
		lf := map[string]interface{}{}
		lf["msg"] = msg.Value
		log.WithFields(lf).Info(fmt.Sprintf("[kafka-publisher] topic: %s,  partition: %d, offset: %d", msg.Topic, partition, offset))
	}
	return nil
}

// NewProducer return message producer
func NewProducer(cfg *Config) Producer {

	m := &producer{}
	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()

	if cfg.Version == "" {
		cfg.Version = defaultVersion
	}

	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		log.Fatal(fmt.Sprintf("parse kafka version got: %v", err))
	}

	config.Producer.Idempotent = cfg.Producer.IdemPotent
	config.Producer.RequiredAcks = sarama.RequiredAcks(cfg.Producer.RequireACK)

	if cfg.Producer.IdemPotent {
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Net.MaxOpenRequests = 1
	}

	config.Version = version

	if len(strings.Trim(cfg.Producer.PartitionStrategy, " ")) == 0 {
		cfg.Producer.PartitionStrategy = "hash"
	}

	strategy, ok := partitiions[cfg.Producer.PartitionStrategy]

	if !ok {
		log.Fatal(fmt.Sprintf("[kafka] invalid producer partition strategy %s", cfg.Producer.PartitionStrategy))
	}

	if cfg.SASL.Enable {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cfg.SASL.User
		config.Net.SASL.Password = cfg.SASL.Password
		config.Net.SASL.Version = cfg.SASL.Version
		config.Net.SASL.Handshake = cfg.SASL.Handshake
		config.Net.SASL.Mechanism = sarama.SASLMechanism(cfg.SASL.Mechanism)
		config.Net.TLS.Enable = true
	}

	// The TLS configuration to use for secure connections if
	// enabled (defaults to nil).
	if config.Net.TLS.Enable || cfg.TLS.Enable {
		config.Net.TLS.Config = createTlsConfig(cfg.TLS)
		config.Net.TLS.Enable = true
	}

	config.Producer.Partitioner = strategy

	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Producer.Timeout = time.Duration(cfg.Producer.TimeoutSecond) * time.Second

	if cfg.Producer.TimeoutSecond < 1 {
		config.Producer.Timeout = defaultTimeout * time.Second
	}

	if cfg.ChannelBufferSize > 0 {
		config.ChannelBufferSize = cfg.ChannelBufferSize
	}

	m.brokers = cfg.Brokers
	m.config = config

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)

	if err != nil {
		log.Fatal(fmt.Sprintf("failed to start Sarama producer:%s", err.Error()))
	}

	m.producer = producer

	return m
}
