package main

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type SinkWriter interface {
	Write(string) error
}

type ConsoleWriter struct{}

func (c *ConsoleWriter) Write(s string) error {
	log.Println("---------->>>", s)
	return nil
}

type KafkaWriter struct {
	Config *Config
	w      *kafka.Writer
}

func (k *KafkaWriter) Write(s string) error {
	k.w.WriteMessages(context.Background(), kafka.Message{Value: []byte(s)})
	return nil
}

func NewSinkWriter(config *Config) SinkWriter {
	if config.Out.Kafka.Hosts != nil {
		w := &kafka.Writer{
			Addr:  kafka.TCP(config.Out.Kafka.Hosts...),
			Topic: config.Out.Kafka.Topic,
		}
		if config.Out.Kafka.Username != "" {
			mechanism, _ := scram.Mechanism(scram.SHA256, config.Out.Kafka.Username, config.Out.Kafka.Password)
			w.Transport = &kafka.Transport{
				SASL: mechanism,
				TLS:  &tls.Config{},
			}
		}

		return &KafkaWriter{
			Config: config,
			w:      w,
		}
	}
	return &ConsoleWriter{}
}
