package entities

type EnvironmentIssuesModels struct {
	ID    uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	Name  string `gorm:"column:name;type:varchar(255)" json:"name"`
	Photo string `gorm:"column:photo;type:varchar(255)" json:"photo"`
}

func (EnvironmentIssuesModels) TableName() string {
	return "environment_issues"
}
