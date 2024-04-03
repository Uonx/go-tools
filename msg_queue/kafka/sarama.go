package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type Sarama struct {
	Address string
	Client  sarama.Client
	Err     error
}

var SaramaMq Sarama
var saramaOnce sync.Once

func NewSaramaMq(address string) {
	saramaOnce.Do(func() {
		SaramaMq = Sarama{
			Address: address,
		}
		SaramaMq.init()
	})
}

func (s *Sarama) init() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Fetch.Max = 10
	c, err := sarama.NewClient(strings.Split(s.Address, ";"), config)
	if err != nil {
		log.Fatalf("sarama init err: %v \n", err)
	}
	s.Client = c
}

func (s *Sarama) SyncProducer(topic, message string) error {
	producer, err := sarama.NewSyncProducerFromClient(s.Client)
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sarama) SyncProducers(m map[string][]string) error {
	producer, err := sarama.NewSyncProducerFromClient(s.Client)
	if err != nil {
		return err
	}
	defer producer.Close()
	var msgs []*sarama.ProducerMessage
	for k, v := range m {
		for _, m := range v {
			msgs = append(msgs, &sarama.ProducerMessage{
				Topic: k,
				Value: sarama.StringEncoder(m),
			})
		}
	}

	return producer.SendMessages(msgs)
}

func (s *Sarama) ConsumerGroup(group_id string, topics []string, handler sarama.ConsumerGroupHandler) error {
	group, err := sarama.NewConsumerGroupFromClient(group_id, s.Client)
	if err != nil {
		return err
	}
	defer group.Close()
	for {
		err := group.Consume(context.Background(), topics, handler)
		if err != nil {
			log.Fatalf("group_id [%s] topics [%v] err: %v \n", group_id, topics, err)
			return err
		}
	}
}

func (s *Sarama) Consumer(topic string) {
	consumer, err := sarama.NewConsumerFromClient(s.Client)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	partitions, _ := consumer.Partitions(topic)
	fmt.Println(topic, partitions)
	for _, partition := range partitions {
		pc, _ := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		go func(pc sarama.PartitionConsumer) {
			defer pc.AsyncClose()
			for msg := range pc.Messages() {
				fmt.Println(msg.Value)
			}
		}(pc)
	}
}
