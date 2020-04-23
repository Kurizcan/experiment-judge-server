package model

// apiSever 的数据库
type AnswerModel struct {
	Id         int    `json:"-" gorm:"column:id;primary_key;"`
	GroupId    int    `json:"group_id" gorm:"column:groupId;"`
	StudentId  int    `json:"student_id" gorm:"column:studentId;"`
	ProblemId  int    `json:"problem_id" gorm:"column:problemId"`
	Status     int    `json:"status" gorm:"column:status"`
	Score      int    `json:"score" gorm:"column:score"`
	Submit     string `json:"submit" gorm:"column:submit"`
	Error      string `json:"error" gorm:"column:error"`
	Correct    bool   `json:"correct" gorm:"column:correct"`
	UpdateTime int64  `json:"update_time" gorm:"column:updateTime"`
}

func (a *AnswerModel) TableName() string {
	return "answer"
}

func (a *AnswerModel) Detail(id int) error {
	return DB.ApiDatabase.Where("Id = ?", id).Error
}

func (a *AnswerModel) Update(id int, data map[string]interface{}) error {
	return DB.ApiDatabase.Table(a.TableName()).Where("id = ?", id).Updates(data).Error
}
