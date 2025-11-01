package models

// TestTable represents a demo table for testing database connectivity
type TestTable struct {
	BaseModel
	Foo string `gorm:"type:varchar(255)" json:"foo"`
	Bar int    `gorm:"type:integer" json:"bar"`
}

// TableName specifies the table name for the TestTable model
func (TestTable) TableName() string {
	return "test_tables"
}