package model

import (
	"experiment-judge-server/config"
	"testing"
)

func TestCreateDataBase(t *testing.T) {

	if err := config.Init("G:\\experiment-judge-server\\conf\\config.yml"); err != nil {
		panic(err)
	}
	// init db
	DB.Init()
	defer DB.Close()

	CreateDataBase("problem_112")
}
