package judge

import (
	"experiment-judge-server/model"
	"experiment-judge-server/util"
	"experiment-judge-server/util/message"
	"experiment-judge-server/util/mysql"
	"github.com/lexkong/log"
)

type ProblemHandler struct {
	Error   chan error
	Success chan struct{}
}

func (p *ProblemHandler) CreateProblem(topicProblemMessage message.TopicProblemMessage) {
	// 1. 收到创建题目的消息后，解析消息
	// 2. 将标准答案等存入数据库
	// 3. 执行元数据脚本，创建用于测试的数据库
	// TODO 加上对 message 的校验，加上失败处理进入失败通道等待再次处理
	dataFile, err := util.StoreFile(topicProblemMessage.DataSource)
	if err != nil {
		p.Error <- err
		util.DelFile(dataFile)
		return
	}
	problem := model.ProblemModel{
		ProblemId:  topicProblemMessage.ProblemId,
		DataSource: dataFile,
		Solution:   topicProblemMessage.Solution,
		OutPut:     topicProblemMessage.OutPut,
		DataBase:   util.GetDataBaseName(topicProblemMessage.ProblemId),
	}

	// 首先创建对应数据库
	err = model.CreateDataBase(problem.DataBase)
	if err != nil {
		p.Error <- err
		util.DelFile(dataFile)
		return
	}
	// 执行元数据脚本
	command := &mysql.ProblemCommand{}
	command.SetCommand(problem.DataBase, dataFile)
	log.Infof("data file:%s", dataFile)
	res, err := mysql.RunCommand(command)
	if err != nil || res != 0 {
		log.Errorf(err, "source data fail")
		p.Error <- err
		util.DelFile(dataFile)
		return
	}

	if err := problem.Create(); err != nil {
		p.Error <- err
		util.DelFile(dataFile)
		return
	}

	p.Success <- struct{}{}
	return
}

func (p *ProblemHandler) Close() {
	close(p.Error)
	close(p.Success)
}
