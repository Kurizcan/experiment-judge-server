package judge

import (
	"experiment-judge-server/util/message"
	"github.com/lexkong/log"
)

func RunProblemHandler(topicProblemMessage message.TopicProblemMessage) {
	handler := &ProblemHandler{
		Error:   make(chan error),
		Success: make(chan struct{}),
	}

	go handler.CreateProblem(topicProblemMessage)
	defer handler.Close()

	select {
	case <-handler.Success:
		log.Infof("source data succeeded")
	case err := <-handler.Error:
		log.Errorf(err, " source data fail")
	}
}

func RunAnswerHandler(answerMessage message.TopicAnswerMessage) {
	handler := &AnswerHandler{
		Error:   make(chan error),
		Success: make(chan struct{}),
	}
	go handler.CheckAnswer(answerMessage)
	defer handler.Close()

	select {
	case <-handler.Success:
		log.Infof("check answer succeeded")
	case err := <-handler.Error:
		log.Errorf(err, "check answer fail")
	}
}
