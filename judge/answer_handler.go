package judge

import (
	"encoding/json"
	"errors"
	"experiment-judge-server/model"
	"experiment-judge-server/util"
	"experiment-judge-server/util/message"
	"experiment-judge-server/util/mysql"
	"experiment-judge-server/util/redis"
	"fmt"
	"github.com/lexkong/log"
	"strings"
)

const (
	AllRight    = 20
	PartlyRight = 10
	WrongAnswer = 5
	HaveError   = 2
)

type AnswerHandler struct {
	Error   chan error
	Success chan struct{}
}

func (a *AnswerHandler) CheckAnswer(topicAnswerMessage message.TopicAnswerMessage) {
	// 1. 在数据库中查找对应的数据库是否存在
	// 2. 进行查询得到结果
	// 3. 与正确结果进行比较
	problem := model.ProblemModel{}
	if err := problem.Detail(topicAnswerMessage.ProblemId); err != nil {
		a.returnError(err)
		return
	}

	// 用于更新 api server
	answer := model.AnswerModel{}
	if err := answer.Detail(topicAnswerMessage.AnswerId); err != nil {
		a.returnError(err)
		return
	}

	// 查看此答案提交的时间是否是旧的
	if answer.UpdateTime > topicAnswerMessage.UpdateTime {
		log.Infof("answer.UpdateTime:%d, topicAnswerMessage.UpdateTime:%d", answer.UpdateTime, topicAnswerMessage.UpdateTime)
		a.returnError(errors.New("the answer already handler"))
		return
	}

	// TODO 对待运行的 sql 进行校验, 先选择数据库，然后运行结果
	sql := topicAnswerMessage.Submit
	command := &mysql.AnswerCommand{}
	command.SetCommand(problem.DataBase, sql)
	_, _ = mysql.RunCommand(command)

	outPutData := model.OutPutData{}
	err := json.Unmarshal(problem.OutPut, &outPutData)
	if err != nil {
		a.returnError(err)
		return
	}
	log.Infof("outputData:%v", outPutData)
	score, errs := check(command, outPutData)

	// TODO 根据得到的结果进行比对
	// 1. 运行有错误信息	- 2
	// 2. 运行无错误信息 - 答案完全正确 20, 部分正确 10， 不正确 5
	// 3. 更新 apiSever 结果
	// 4. 删除对应 runId 缓存
	status, correct := util.ProblemSubmitStatus[util.WRONG], false
	if score == AllRight {
		status = util.ProblemSubmitStatus[util.ACCEPT]
		correct = true
	}

	err = answer.Update(topicAnswerMessage.AnswerId, map[string]interface{}{
		"Score":   score,
		"Error":   errs,
		"Status":  status,
		"Correct": correct,
	})

	if err != nil {
		a.returnError(err)
		return
	}

	// 万一缓存删除失败了，由过期时间进行淘汰
	runIdKeys := redis.GetRunIdStatusKey(topicAnswerMessage.AnswerId, answer.StudentId)
	experimentKeys := redis.GetGroupSidKey(answer.GroupId, answer.StudentId)
	err = redis.Client.Del(runIdKeys)
	if err != nil {
		log.Errorf(err, "redis del is fail, runIdKeys key: %s", runIdKeys)
	}
	err = redis.Client.HDel(experimentKeys, fmt.Sprintf("%d", topicAnswerMessage.ProblemId))
	if err != nil {
		log.Errorf(err, "redis del is fail, experimentKeys key: %s", experimentKeys)
	}
	a.returnSuccess()
	return
}

func (a *AnswerHandler) returnError(err error) {
	a.Error <- err
}

func (a *AnswerHandler) returnSuccess() {
	a.Success <- struct{}{}
}

func (a *AnswerHandler) Close() {
	close(a.Error)
	close(a.Success)
}

func check(command *mysql.AnswerCommand, output model.OutPutData) (int, string) {
	rightBool := true
	wrongBool := false
	result := command.Output.String()
	enomsg := command.Enomsg.String()
	// 先检查运行信息是否有误
	enomsgSlice := strings.Split(enomsg, "\n")
	log.Infof("len of es: %d", len(enomsgSlice)) // len > 2 就有错误，默认会输出一条 warning 和 一空行
	if len(enomsgSlice) >= 2 {
		return HaveError, strings.Join(enomsgSlice[1:len(enomsgSlice)-1], "\n")
	}

	// 运行信息无错，检查结果，首先按照换行符进行分隔结果
	resultSlice := strings.Split(result, "\n")
	log.Infof("result:%s,\n output:%v\n", result, output)
	for i := 0; i < len(resultSlice)-1; i++ {
		// 对于每一行，按照空格分隔字符串
		// 最后有一行是空行，默认不处理
		line := strings.Fields(resultSlice[i])
		for j, item := range line {
			if i == 0 {
				// headers compare
				if j < len(output.Headers) && item == output.Headers[j] {
					wrongBool = true
					log.Infof("this line %s is true with header %s", item, output.Headers[j])
				} else {
					rightBool = false
				}
			} else {
				// rows compare
				if i-1 < len(output.Rows) && j < len(output.Rows[i-1]) && item == output.Rows[i-1][j] {
					wrongBool = true
					log.Infof("this line %s is true with row %s", item, output.Rows[i-1][j])
				} else {
					rightBool = false
				}
			}
		}
	}

	if rightBool {
		// 全对
		return AllRight, ""
	} else if !wrongBool {
		// 全错
		return WrongAnswer, "your output were all wrong"
	} else {
		return PartlyRight, "your output were Partly wrong"
	}
}
