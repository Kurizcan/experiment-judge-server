package model

type ProblemModel struct {
	Id         int    `json:"-" gorm:"column:id;primary_key;"`
	ProblemId  int    `json:"problem_id" gorm:"column:problemId;"`
	DataSource string `json:"data_source" gorm:"column:dataSource;"`
	Solution   string `json:"solution" gorm:"column:solution;"`
	OutPut     []byte `json:"out_put" gorm:"column:output;"`
	DataBase   string `json:"data_base" gorm:"column:dataBase;"`
}

type OutPutData struct {
	Headers []string   `json:"headers"`
	Rows    []RowsData `json:"rows"`
}

type RowsData []string

func (p *ProblemModel) TableName() string {
	return "problem"
}

func (p *ProblemModel) Create() error {
	return DB.Self.Create(&p).Error
}

func (p *ProblemModel) Detail(problemId int) error {
	return DB.Self.Where("problemId = ?", problemId).Find(&p).Error
}
