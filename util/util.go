package util

import (
	"fmt"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"github.com/teris-io/shortid"
	"os"
)

const (
	Student   = 0
	Teacher   = 1
	Admin     = 2
	NEW       = "new"
	NOSUBMIT  = "no_submit"
	SUBMIT    = "submit"
	COMPLETED = "completed"
	RUNNING   = "running"
	ACCEPT    = "accept"
	WRONG     = "wrong"
	EMPTY     = -1
)

var ExperimentStudentStatus = map[string]int{
	NEW:       0,
	NOSUBMIT:  1,
	SUBMIT:    2,
	COMPLETED: 3,
}

var ProblemSubmitStatus = map[string]int{
	RUNNING: 1,
	ACCEPT:  2,
	WRONG:   3,
}

func GenShortId() (string, error) {
	return shortid.Generate()
}

func getDataScoreFileName() string {
	name, _ := GenShortId()
	root, _ := os.Getwd()
	return fmt.Sprintf("%s\\%s\\%s.sql", root, viper.GetString("data_scour"), name)
}

func StoreFile(data []byte) (string, error) {
	fileName := getDataScoreFileName()
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("open file fail", err)
		return "", err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		log.Fatal("write data fail", err)
		return "", nil
	}
	return fileName, nil
}

func DelFile(fileName string) {
	_ = os.Remove(fileName)
}

func GetDataBaseName(id int) string {
	return fmt.Sprintf("problem_%d", id)
}
