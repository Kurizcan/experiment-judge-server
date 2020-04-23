package message

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type KafkaClient struct {
	Client          sarama.Client
	Config          *sarama.Config
	AnswerConsumer  sarama.ConsumerGroup
	ProblemConsumer sarama.ConsumerGroup
}

var singleton *KafkaClient
var once sync.Once

func GetKafkaClient() *KafkaClient {
	once.Do(func() {
		// Init config, specify appropriate version
		config := sarama.NewConfig()
		config.Version = sarama.V1_0_0_0
		config.Consumer.Return.Errors = true
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
		// Start with a client
		client, err := sarama.NewClient(viper.GetStringSlice("mq.kafka.brokers"), config)
		if err != nil {
			panic(err)
		}
		// Start a new consumer group
		answerGroup, err := sarama.NewConsumerGroupFromClient(viper.GetString("mq.consumer.answer"), client)
		if err != nil {
			panic(err)
		}
		log.Info("answerGroup is create")
		problemGroup, err := sarama.NewConsumerGroupFromClient(viper.GetString("mq.consumer.problem"), client)
		if err != nil {
			panic(err)
		}
		log.Info("problemGroup is create")
		singleton = &KafkaClient{
			Client:          client,
			Config:          config,
			AnswerConsumer:  answerGroup,
			ProblemConsumer: problemGroup,
		}
	})
	return singleton
}

func (k *KafkaClient) Consumer() {
	go runConsumer(k.ProblemConsumer, problemConsumerGroupHandler{}, []string{TopicProblem})
	go runConsumer(k.AnswerConsumer, answerConsumerGroupHandler{}, []string{TopicAnswer})
}

func (k *KafkaClient) Close() {
	_ = k.ProblemConsumer.Close()
	_ = k.AnswerConsumer.Close()
	_ = k.Client.Close()
}

func runConsumer(consumer sarama.ConsumerGroup, handler sarama.ConsumerGroupHandler, topics []string) {
	// Track errors
	go func() {
		for err := range consumer.Errors() {
			log.Errorf(err, " consume answer message failed")
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		err := consumer.Consume(ctx, topics, handler)
		if err != nil {
			log.Errorf(err, " consume answer message failed")
		}
	}
}
