package mysql

import (
	"bytes"
	"fmt"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"os/exec"
	"syscall"
)

type Command interface {
	Run() (uint32, error)
}

func RunCommand(command Command) (uint32, error) {
	return command.Run()
}

func setMysqlCommand() string {
	return "mysql"
}

func setArgs(dataBase, source string) []string {
	arg := make([]string, 0)
	arg = append(arg, fmt.Sprintf("--host=%s", viper.GetString("db.host")))
	arg = append(arg, fmt.Sprintf("--user=%s", viper.GetString("db.username")))
	arg = append(arg, fmt.Sprintf("--password=%s", viper.GetString("db.password")))
	arg = append(arg, fmt.Sprintf("--database=%s", dataBase))
	arg = append(arg, fmt.Sprintf("--execute=%s", source))
	return arg
}

type ProblemCommand struct {
	cmd *exec.Cmd
}

func (p *ProblemCommand) Run() (uint32, error) {
	err := p.cmd.Run()
	if err != nil {
		log.Errorf(err, "ProblemCommand run fail, the command %s %s", p.cmd.Path, p.cmd.Args)
		return 0, err
	}
	res := p.cmd.ProcessState.Sys().(syscall.WaitStatus).ExitCode
	log.Infof("the command %s %s res:%d", p.cmd.Path, p.cmd.Args, res)
	return res, nil
}

// mysql -u root --password=experiment problem_101 < G:\experiment-judge-server\data\78eENa3ZR.sql
func (p *ProblemCommand) SetCommand(dataBase, dataSource string) {
	args, cmd := setArgs(dataBase, fmt.Sprintf("source %s", dataSource)), setMysqlCommand()
	command := exec.Command(cmd, args...)
	p.cmd = command
}

type AnswerCommand struct {
	cmd    *exec.Cmd
	Output bytes.Buffer
	Enomsg bytes.Buffer
}

func (a *AnswerCommand) SetCommand(dataBase, sql string) {
	args, cmd := setArgs(dataBase, sql), setMysqlCommand()
	command := exec.Command(cmd, args...)
	a.cmd = command
	a.cmd.Stdout = &a.Output
	a.cmd.Stderr = &a.Enomsg
}

func (a *AnswerCommand) Run() (uint32, error) {
	err := a.cmd.Run()
	if err != nil {
		log.Errorf(err, "AnswerCommand run fail, the AnswerCommand %s %s", a.cmd.Path, a.cmd.Args)
		return 0, err
	}
	res := a.cmd.ProcessState.Sys().(syscall.WaitStatus).ExitCode
	log.Infof("the command %s %s res:%d", a.cmd.Path, a.cmd.Args, res)
	return res, nil
}
