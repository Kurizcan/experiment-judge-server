package util

import (
	"bytes"
	"experiment-judge-server/config"
	"experiment-judge-server/model"
	"fmt"
	"os/exec"
	"syscall"
	"testing"
)

func TestGetDataBaseName(t *testing.T) {
	// init config
	if err := config.Init("G:\\experiment-judge-server\\conf\\config.yml"); err != nil {
		panic(err)
	}
	model.DB.Init()
	defer model.DB.Close()
	//fmt.Println(getDataScoreFileName())
	// mysql -h 127.0.0.1 -u root --password=123456 problem_118 < G:\\experiment-judge-server\\data\\cM9HF-qZR.sql
	res := "mysql"
	arg := make([]string, 0)
	arg = append(arg, "--host=localhost")
	arg = append(arg, "--user=root")
	arg = append(arg, "--password=123456")
	arg = append(arg, "--database=problem_112")
	arg = append(arg, "--execute=source G:\\experiment-judge-server\\data\\cM9HF-qZR.sql")
	command := exec.Command(res, arg...)
	var output bytes.Buffer
	var enomsg bytes.Buffer
	command.Stdout = &output
	command.Stderr = &enomsg
	err := command.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(command.ProcessState.Sys().(syscall.WaitStatus).ExitCode)
	fmt.Println(command.Path)
	fmt.Println(command.Args)
	fmt.Println(output.String())
	fmt.Println(enomsg.String())
}
