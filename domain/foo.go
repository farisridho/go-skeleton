package domain

type Foo struct {
	ID   string `gorm:"primary_key;type:uuid;" json:"Id"`
	Name string `gorm:"-" json:"Name"`
}

func (t Foo) TableName() string {
	return "foo"
}
