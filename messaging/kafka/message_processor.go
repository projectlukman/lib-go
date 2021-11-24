// Package kafka messaging
// @author Daud Valentino
package kafka

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type MessageProcessorFunc func(*MessageDecoder)

// MessageProcessor contract message consumer processor
type MessageProcessor interface {
	Processor(decoder *MessageDecoder) error
}

// MessageDecoder decoder message data  on topic
type MessageDecoder struct {
	Body      []byte
	Key       []byte
	Topic     string
	Partition int32
	TimeStamp time.Time
	Offset    int64
	Commit    func(*MessageDecoder)
}

// DecodeJSON decode kafka message byte to struct
func (decoder *MessageDecoder) DecodeJSON(out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("%s", "output destination cannot addressable")
	}

	return json.Unmarshal(decoder.Body, out)
}

// MessageEncoder message encoder  publish message to kafka
type MessageEncoder interface {
	Encode() ([]byte, error)
	Key() string
	Length() int
}
