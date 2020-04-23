package message

// 为了解决 judge 与 message 包循环依赖的问题，将这两个消息体抽到了 util 下

type TopicProblemMessage struct {
	ProblemId  int
	DataSource []byte
	Solution   string
	OutPut     []byte
}

type TopicAnswerMessage struct {
	AnswerId   int
	ProblemId  int
	Submit     string
	UpdateTime int64
}
