package message

import (
	"encoding/json"
	"experiment-judge-server/judge"
	"experiment-judge-server/util/message"
	"github.com/Shopify/sarama"
	"github.com/lexkong/log"
)

type MQ interface {
	Consumer()
	Close()
}

const (
	TopicProblem = "problem_experiment" // 2 个副本， 3 个分区
	TopicAnswer  = "answer_experiment"  // 2 个副本， 3 个分区
)

type answerConsumerGroupHandler struct{}

func (answerConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (answerConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h answerConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Infof("Message topic:%q partition:%d offset:%d value:%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		answerMessage := message.TopicAnswerMessage{}
		err := json.Unmarshal(msg.Value, &answerMessage)
		if err != nil {
			log.Error("json fail ", err)
		}
		log.Infof("received msg:%v", answerMessage)

		// TODO 将消息传递给 answer_judge 可用 channel 等
		go judge.RunAnswerHandler(answerMessage)

		sess.MarkMessage(msg, "")
	}
	return nil
}

type problemConsumerGroupHandler struct{}

func (problemConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (problemConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h problemConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Infof("Message topic:%q partition:%d offset:%d value:%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		problemMessage := message.TopicProblemMessage{}
		err := json.Unmarshal(msg.Value, &problemMessage)
		if err != nil {
			log.Error("json fail ", err)
		}
		log.Infof("received TopicProblemMessage, pid:%d, solution:%s\n", problemMessage.ProblemId, problemMessage.Solution)

		go judge.RunProblemHandler(problemMessage)

		// TODO 将消息传递给 problem_handler 可用 channel 等

		sess.MarkMessage(msg, "")
	}
	return nil
}
