package mysql

import (
	"bytes"
	"encoding/json"
	"experiment-judge-server/config"
	"experiment-judge-server/model"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

func TestAnswerCommand_Run(t *testing.T) {

	if err := config.Init("G:\\experiment-judge-server\\conf\\config.yml"); err != nil {
		panic(err)
	}
	// init db
	model.DB.Init()
	defer model.DB.Close()

	res := "mysql"
	arg := make([]string, 0)
	arg = append(arg, "--host=localhost")
	arg = append(arg, "--user=root")
	arg = append(arg, "--password=123456")
	arg = append(arg, "--database=problem_112")
	arg = append(arg, "--execute=select * from A where A.id = 15")
	//arg = append(arg, "--execute=select customers.name as 'Customers' from customers where customers.id not in (select customerid from orders);")
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
	fmt.Println("err:", enomsg.String())

	es := strings.Split(enomsg.String(), "\n")
	fmt.Println("len of es:", len(es)) // len > 2 就有错误

	fmt.Println("es[2]:", es[1])

	problemModel := model.ProblemModel{}
	err = problemModel.Detail(119)
	if err != nil {
		fmt.Println("fail")
	}
	outPutData := model.OutPutData{}
	json.Unmarshal(problemModel.OutPut, &outPutData)
	fmt.Println(outPutData)

	// 首先按照换行符进行分隔结果
	s := strings.Split(output.String(), "\n")
	fmt.Println("s len:", len(s))
	for i := 0; i < len(s)-1; i++ {
		// 对于每一行，按照空格分隔字符串
		// 最后有一行是空行，默认不处理
		fmt.Println(s[i])
		line := strings.Fields(s[i])
		for j, item := range line {
			if i == 0 {
				// headers compare
				fmt.Print(item)
				if j < len(outPutData.Headers) && item == outPutData.Headers[j] {
					fmt.Println(" is true, header")
				} else {
					fmt.Println(" is false, header")
				}
			} else {
				// rows compare
				fmt.Print(item)
				if i-1 < len(outPutData.Rows) && j < len(outPutData.Rows[i-1]) && item == outPutData.Rows[i-1][j] {
					fmt.Println(" is true, row")
				} else {
					fmt.Println(" is false, row")
				}
			}
		}
	}

}
