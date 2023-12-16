package entities

type StatusSeederModels struct {
	ID         uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	IsExecuted bool   `gorm:"column:is_executed;default:false" json:"is_executed"`
}

func (StatusSeederModels) TableName() string {
	return "status_seeder"
}
